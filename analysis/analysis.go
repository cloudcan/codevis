package analysis

import (
	"errors"
	"fmt"
	"github.com/cloudcan/codevis/graphdb"
	"go/token"
	"go/types"
	callgraph2 "golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"log"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"time"
)

type Config struct {
	Dir string `json:"dir"`
}

// analysis result
type Result struct {
	prog      *ssa.Program
	callgraph *callgraph2.Graph
}
type NodeHandler func(Node) error
type EdgeHandler func(Edge) error

// static code analysis
func Analysis(config Config) (r *Result, err error) {
	// load packages
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  config.Dir,
	}
	initial, err := packages.Load(cfg)
	if err != nil {
		err = errors.New(fmt.Sprintf("load packages error,cause:%s", err))
		return
	}
	if packages.PrintErrors(initial) > 0 {
		err = errors.New("load packages  include error")
		return
	}
	//  build ssa program
	prog, pkgs := ssautil.AllPackages(initial, 0)
	prog.Build()
	// get call graph
	mains := ssautil.MainPackages(pkgs)
	pCfg := &pointer.Config{
		Mains:          mains,
		Reflection:     true,
		BuildCallGraph: true,
	}
	result, err := pointer.Analyze(pCfg)
	if err != nil {
		err = errors.New(fmt.Sprintf("pointer analysis err,cause:%s", err))
		return
	}
	return &Result{
		prog:      prog,
		callgraph: result.CallGraph,
	}, nil
}

// refine analysis result
func (r *Result) Refine() (g *Graph, err error) {
	var (
		files    = make(map[string]*File, 0)
		pkgs     = make(map[string]*Package, 0)
		elements = make([]Element, 0)
		nodes    = make([]Node, 0)
		edges    = make([]Edge, 0)
	)
	// resolve files
	fset := r.prog.Fset
	fset.Iterate(func(file *token.File) bool {
		path := file.Name()
		files[path] = NewFile(filepath.Base(path), path, "", file.LineCount())
		return file != nil
	})
	// resolve packages
	allPkg := r.prog.AllPackages()
	if len(allPkg) > 0 {
		for _, pkg := range allPkg {
			// handle package
			ppkg := NewPackage(pkg.Pkg.Name(), pkg.Pkg.Path())
			pkgs[ppkg.path] = ppkg
			// handle pkg contains files
			scope := pkg.Pkg.Scope()
			numChildren := scope.NumChildren()
			for i := 0; i < numChildren; i++ {
				value := reflect.ValueOf(*scope.Child(i))
				comment := value.FieldByName("comment").String()
				file := files[comment]
				file.pkg = ppkg.path
			}
			// handle file imports
			// handle elements
			for _, member := range pkg.Members {
				var (
					ele      Element
					position token.Position
				)
				switch m := member.(type) {
				case *ssa.Function:
					position = fset.Position(m.Pos())
					ele = NewFunction(m.Name(), m.Pkg.Pkg.Path(), m.Signature.String(), position.Filename, Position{
						column: position.Column,
						line:   position.Line,
					})
				case *ssa.Global:
					position = fset.Position(m.Pos())
					ele = NewGlobal(m.Name(), m.Pkg.Pkg.Path(), m.Type().String(), position.Filename, Position{
						column: position.Column,
						line:   position.Line,
					})
				case *ssa.NamedConst:
					position = fset.Position(m.Pos())
					ele = NewNamedConst(m.Name(), m.Package().Pkg.Path(), position.Filename, m.Type().String(), m.Value.Value.String(), Position{
						column: position.Column,
						line:   position.Line,
					})
				case *ssa.Type:
					position = fset.Position(m.Pos())
					underlying := m.Object().Type().Underlying()
					var (
						fields  []string
						methods []string
					)
					if value, ok := underlying.(*types.Struct); ok {
						count := value.NumFields()
						fields = make([]string, count)
						for i := 0; i < count; i++ {
							fields[i] = fmt.Sprintf("%s:%s", value.Field(i).Name(), value.Field(i).Type())
						}
					}
					ele = NewType(m.Name(), m.Package().Pkg.Name(), position.Filename, underlying.String(), fields, methods, Position{
						column: position.Column,
						line:   position.Line,
					})
				default:
					log.Print("unknow type:", m)
				}
				elements = append(elements, ele)
			}
		}
	}
	// add elements
	for _, ele := range elements {
		if f := files[ele.File()]; f != nil {
			edges = append(edges, &Declare{
				from: f,
				to:   ele,
			})
		}
		if p := pkgs[ele.Pkg()]; p != nil {
			edges = append(edges, &Contains{
				from: p,
				to:   ele,
			})
		}
		nodes = append(nodes, ele)
	}
	// add file
	for _, file := range files {
		if p := pkgs[file.pkg]; p != nil {
			edges = append(edges, &Contains{
				from: p,
				to:   file,
			})
		}
		nodes = append(nodes, file)
	}
	// add package
	for _, pkg := range pkgs {
		nodes = append(nodes, pkg)
	}
	// add root relationship
	root := NewProgram()
	for _, node := range nodes {
		edges = append(edges, &Belong{
			root: root,
			node: node,
		})
	}
	nodes = append(nodes, root)
	return &Graph{
		Root:  root,
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// save graph to db
func (g *Graph) Persistence() {
	// for stat
	var (
		start  = time.Now()
		ticker = time.NewTicker(1 * time.Second)
		count  = uint32(0)
		stopCh = make(chan struct{})
	)
	go func() {
		prev := uint32(0)
		for {
			select {
			case <-ticker.C:
				now := atomic.LoadUint32(&count)
				log.Printf("save speed:%d/s\n", now-prev)
				prev = now
			case <-stopCh:
				ticker.Stop()
				log.Printf("save complete,total:%d nodes and %d edges ,cost:%s", len(g.Nodes), len(g.Edges), time.Now().Sub(start))
				return
			}
		}

	}()
	// update deprecated node and edge
	_, err := graphdb.Exec(fmt.Sprintf("match (n:%s{name:'%s'})--(m)  detach delete n,m", lProg, g.Root.name), nil)
	if err != nil {
		log.Print("update deprecated node error,cause:", err)
	}
	// save node
	if len(g.Nodes) > 0 {
		log.Print("start save node ...")
		for _, node := range g.Nodes {
			atomic.AddUint32(&count, 1)
			_, err := graphdb.Exec(fmt.Sprintf("create(:%s%s)", node.Label(), node.Body()), nil)
			if err != nil {
				log.Print("save node err,cause:", err)
			}
		}
	}
	// create index
	graphdb.CreateIndex(string(lProg), "name", "uuid")
	graphdb.CreateIndex(string(lPackage), "name", "path", "uuid")
	graphdb.CreateIndex(string(lFile), "name", "path", "pkg", "uuid")
	graphdb.CreateIndex(string(lGlobal), "name", "file", "pkg", "uuid")
	graphdb.CreateIndex(string(lConst), "name", "file", "pkg", "uuid")
	graphdb.CreateIndex(string(lFunc), "name", "file", "pkg", "uuid")
	graphdb.CreateIndex(string(lType), "name", "file", "pkg", "uuid")
	// save edge
	if len(g.Edges) > 0 {
		log.Print("start save edge ...")
		for _, edge := range g.Edges {
			atomic.AddUint32(&count, 1)
			_, err := graphdb.Exec(fmt.Sprintf("match (n:%s{uuid:'%s'}),(m:%s{uuid:'%s'})  create (n)-[:%s]->(m)", edge.From().Label(), edge.From().Id(), edge.To().Label(), edge.To().Id(), edge.RelationshipType()), nil)
			if err != nil {
				log.Print("save edge err,cause:", err)
			}
		}
	}
	close(stopCh)
}
