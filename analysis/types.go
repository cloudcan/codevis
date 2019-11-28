package analysis

import (
	"fmt"
	"github.com/cloudcan/codevis/graphdb"
	"github.com/google/uuid"
	"os"
	"path"
	"strings"
)

type NodeLabel string

const (
	lProg    NodeLabel = "Program"
	lPackage NodeLabel = "Package"
	lFile    NodeLabel = "File"
	lFunc    NodeLabel = "Function"
	lMethod  NodeLabel = "Method"
	lGlobal  NodeLabel = "Global"
	lConst   NodeLabel = "Const"
	lType    NodeLabel = "Type"
)

type RType string

const (
	rContains    RType = "Contains"
	rDeclare     RType = "Declare"
	rCall        RType = "Call"
	rImport      RType = "Import"
	rRecv        RType = "Receive"
	rImpletments RType = "Implements"
	rBelong      RType = "Belong"
)

// node
type Node interface {
	Label() NodeLabel
	Body() string
	Id() string
	fmt.Stringer
}
type Program struct {
	name string
	uuid string
}

func NewProgram() *Program {
	dir, _ := os.Getwd()
	name := path.Base(dir)
	return &Program{name: name, uuid: uuid.New().String()}
}

func (p *Program) Id() string {
	return p.uuid
}
func (p *Program) Label() NodeLabel {
	return lProg
}
func (p *Program) Body() string {
	return fmt.Sprintf("{name:'%s',uuid:'%s'}", p.name, p.uuid)
}
func (p *Program) String() string {
	return fmt.Sprintf("(%s:%s)", p.Label(), p.Body())
}

// package node
type Package struct {
	name string
	path string
	uuid string
}

func NewPackage(name string, path string) *Package {
	return &Package{name: name, path: path, uuid: uuid.New().String()}
}

func (p *Package) Id() string {
	return p.uuid
}
func (p *Package) Label() NodeLabel {
	return lPackage
}
func (p *Package) Body() string {
	return fmt.Sprintf("{name:'%s',path:'%s',uuid:'%s'}", p.name, p.path, p.Id())
}
func (p *Package) String() string {
	return fmt.Sprintf("(%s:%s)", p.Label(), p.Body())
}

type Position struct {
	column int
	line   int
}

func (p *Position) Pos() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

// element node
type Element interface {
	Node
	Pos() string
	Pkg() string
	File() string
}

// file node
type File struct {
	name  string
	path  string
	pkg   string
	lines int
	uuid  string
}

func NewFile(name string, path string, pkg string, lines int) *File {
	return &File{name: name, path: path, pkg: pkg, lines: lines, uuid: uuid.New().String()}
}

func (f *File) Id() string {
	return f.uuid
}
func (f *File) Label() NodeLabel {
	return lFile
}
func (f *File) Body() string {
	return fmt.Sprintf("{name:'%s',path:'%s',pkg:'%s',lines:'%d',uuid:'%s'}", f.name, f.path, f.pkg, f.lines, f.Id())
}
func (f *File) String() string {
	return fmt.Sprintf("(%s:%s)", f.Label(), f.Body())
}

// function node
type Function struct {
	name string
	pkg  string
	sign string
	file string
	uuid string
	Position
}

func NewFunction(name string, pkg string, sign string, file string, position Position) *Function {
	return &Function{name: name, pkg: pkg, sign: sign, file: file, Position: position, uuid: uuid.New().String()}
}
func (f *Function) Id() string {
	return f.uuid
}
func (f *Function) Pkg() string {
	return f.pkg
}
func (f *Function) File() string {
	return f.file
}
func (f *Function) Label() NodeLabel {
	return lFunc
}
func (f *Function) Body() string {
	return fmt.Sprintf("{name:'%s',pkg:'%s',sign:'%s',file:'%s',pos:'%s',,uuid:'%s'}", f.name, f.pkg, f.sign, f.file, f.Pos(), f.Id())
}
func (f *Function) String() string {
	return fmt.Sprintf("(%s%s)", f.Label(), f.Body())
}

// global node
type Global struct {
	name string
	pkg  string
	typ  string
	file string
	uuid string
	Position
}

func NewGlobal(name string, pkg string, typ string, file string, position Position) *Global {
	return &Global{name: name, pkg: pkg, typ: typ, file: file, Position: position, uuid: uuid.New().String()}
}
func (g *Global) Id() string {
	return g.uuid
}
func (g *Global) Pkg() string {
	return g.pkg
}
func (g *Global) File() string {
	return g.file
}
func (g *Global) Label() NodeLabel {
	return lGlobal
}
func (g *Global) Body() string {
	escapePath := strings.ReplaceAll(g.file, "/", "//")
	return fmt.Sprintf("{name:'%s',pkg:'%s',typ:'%s',file:'%s',pos:'%s',uuid:'%s'}", g.name, g.pkg, g.typ, escapePath, g.Pos(), g.Id())
}
func (g *Global) String() string {
	return fmt.Sprintf("(%s%s)", g.Label(), g.Body())
}

// named const node
type NamedConst struct {
	name  string
	pkg   string
	file  string
	typ   string
	value string
	uuid  string
	Position
}

func NewNamedConst(name string, pkg string, file string, typ string, value string, position Position) *NamedConst {
	return &NamedConst{name: name, pkg: pkg, file: file, typ: typ, value: value, Position: position, uuid: uuid.New().String()}
}

func (n *NamedConst) Id() string {
	return n.uuid
}
func (n *NamedConst) Pkg() string {
	return n.pkg
}
func (n *NamedConst) File() string {
	return n.file
}
func (n *NamedConst) Label() NodeLabel {
	return lConst
}
func (n *NamedConst) Body() string {
	return fmt.Sprintf("{name:'%s',pkg:'%s',typ:'%s',file:'%s',value:'%s',pos:'%s',uuid:'%s'}", n.name, n.pkg, n.typ, n.file, graphdb.Escape(n.value), n.Pos(), n.Id())
}
func (n *NamedConst) String() string {
	return fmt.Sprintf("(%s%s)", n.Label(), n.Body())
}

// type node
type Type struct {
	name       string
	pkg        string
	file       string
	underlying string
	fields     []string
	methods    []string
	uuid       string
	Position
}

func NewType(name string, pkg string, file string, underlying string, fields []string, methods []string, position Position) *Type {
	return &Type{name: name, pkg: pkg, file: file, underlying: underlying, fields: fields, methods: methods, Position: position, uuid: uuid.New().String()}
}
func (t *Type) Id() string {
	return t.uuid
}
func (t *Type) Pkg() string {
	return t.pkg
}
func (t *Type) File() string {
	return t.file
}
func (t *Type) Label() NodeLabel {
	return lType
}
func (t *Type) Body() string {
	return fmt.Sprintf("{name:'%s',pkg:'%s',underlying:'%s',file:'%s',fields:'%s',methods:'%s',pos:'%s',uuid:'%s'}", t.name, t.pkg, t.underlying, t.file, t.fields, t.methods, t.Pos(), t.Id())
}
func (t *Type) String() string {
	return fmt.Sprintf("(%s%s)", t.Label(), t.Body())
}

type Edge interface {
	From() string
	To() string
	RelationshipType() RType
	fmt.Stringer
}

// import relationship
type Import struct {
	from *File
	to   *Package
}

func (i *Import) RelationshipType() RType {
	return rImport
}
func (i *Import) From() string {
	return i.from.Id()
}

func (i *Import) To() string {
	return i.to.Id()
}
func (i *Import) String() string {
	return fmt.Sprintf("%s-Import->%s", i.From(), i.To())
}

// contain relationship
type Contains struct {
	from *Package
	to   Node
}

func (c *Contains) RelationshipType() RType {
	return rContains
}
func (c *Contains) From() string {
	return c.from.Id()
}

func (c *Contains) To() string {
	return c.to.Id()
}
func (c *Contains) String() string {
	return fmt.Sprintf("%s-Contains->%s", c.From(), c.To())
}

// declare relationship
type Declare struct {
	from *File
	to   Element
}

func (d *Declare) RelationshipType() RType {
	return rDeclare
}
func (d *Declare) From() string {
	return d.from.Id()
}

func (d *Declare) To() string {
	return d.to.Id()
}
func (d *Declare) String() string {
	return fmt.Sprintf("%s-Declare->%s", d.From(), d.To())
}

// function call relationship
type Call struct {
	caller *Function
	callee *Function
}

func (c *Call) RelationshipType() RType {
	return rCall
}
func (c *Call) From() string {
	return c.caller.Id()
}

func (c *Call) To() string {
	return c.callee.Id()
}
func (c *Call) String() string {
	return fmt.Sprintf("%s-Call->%s", c.From(), c.To())
}

type Receive struct {
	sender   *Type
	receiver *Function
}

func (r *Receive) RelationshipType() RType {
	return rRecv
}
func (r *Receive) From() string {
	return r.sender.Id()
}

func (r *Receive) To() string {
	return r.receiver.Id()
}

func (r *Receive) String() string {
	return fmt.Sprintf("%s<-Receive-%s", r.From(), r.To())
}

type Belong struct {
	root *Program
	node Node
}

func (b *Belong) RelationshipType() RType {
	return rBelong
}
func (b *Belong) From() string {
	return b.root.Id()
}

func (b *Belong) To() string {
	return b.node.Id()
}

func (b *Belong) String() string {
	return fmt.Sprintf("%s-Belong->%s", b.From(), b.To())
}

type Graph struct {
	Root  *Program
	Nodes []Node
	Edges []Edge
}
