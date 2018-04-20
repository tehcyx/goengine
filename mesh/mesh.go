package mesh

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tehcyx/goengine/obj"
	"github.com/tehcyx/goengine/util"
)

const (
	POSITION_VB int = 0
	TEXCOORD_VB int = 1
	NORMAL_VB   int = 2
	INDEX_VB    int = 3
	NUM_BUFFERS int = 4
)

type Vertex struct {
	position mgl32.Vec3
	texCoord mgl32.Vec2
	normal   mgl32.Vec3
}

type Mesh struct {
	vao        uint32              // vertex array object
	vbo        [NUM_BUFFERS]uint32 // vertex buffer object
	model      *obj.IndexedModel
	quickmodel *obj.QuickObjModel
}

func NewMesh(path string) *Mesh {
	m := new(Mesh)
	m.quickInit(obj.LoadObj(path))
	return m
}

// NewMeshFromFile creates a new Mesh from specified string path
func NewMeshFromFile(file string) *Mesh {
	defer util.TimeTrack(time.Now(), "NewMeshFromFile")
	m := new(Mesh)
	// objectModel := obj.NewObjModelFromFile(file)
	// indexedModel := objectModel.ToIndexedModel()
	// m.init(indexedModel)
	m.init(obj.NewObjModelFromFile(file).ToIndexedModel())
	return m
}

func (m *Mesh) create(vertices []Vertex, indices []int) {
	defer util.TimeTrack(time.Now(), "Mesh create -> init")
	model := new(obj.IndexedModel)

	for i := 0; i < len(vertices); i++ {
		model.Positions = append(model.Positions, vertices[i].position)
		model.TexCoords = append(model.TexCoords, vertices[i].texCoord)
		model.Normals = append(model.Normals, vertices[i].normal)
	}

	for i := 0; i < len(indices); i++ {
		model.Indices = append(model.Indices, indices[i])
	}

	m.init(model)

}

func (m *Mesh) init(model *obj.IndexedModel) {
	defer util.TimeTrack(time.Now(), "Mesh init")
	m.model = model

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(int32(NUM_BUFFERS), &m.vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[POSITION_VB])
	gl.BufferData(gl.ARRAY_BUFFER, len(m.model.Positions)*3*4, gl.Ptr(m.model.Positions), gl.STATIC_DRAW) // *4 because  float32 is 4 bytes
	gl.EnableVertexAttribArray(uint32(POSITION_VB))
	gl.VertexAttribPointer(uint32(POSITION_VB), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[TEXCOORD_VB])
	gl.BufferData(gl.ARRAY_BUFFER, len(m.model.TexCoords)*2*4, gl.Ptr(m.model.TexCoords), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(TEXCOORD_VB))
	gl.VertexAttribPointer(uint32(TEXCOORD_VB), 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[NORMAL_VB])
	gl.BufferData(gl.ARRAY_BUFFER, len(m.model.Normals)*3*4, gl.Ptr(m.model.Normals), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(NORMAL_VB))
	gl.VertexAttribPointer(uint32(NORMAL_VB), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.vbo[INDEX_VB])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.model.Indices)*3*4, gl.Ptr(m.model.Indices), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
}

func (m *Mesh) quickInit(model *obj.QuickObjModel) {
	defer util.TimeTrack(time.Now(), "Mesh quickInit")
	m.quickmodel = model

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(int32(NUM_BUFFERS), &m.vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[POSITION_VB])
	gl.BufferData(gl.ARRAY_BUFFER, len(m.quickmodel.Vertices)*3*4, gl.Ptr(m.quickmodel.Vertices), gl.STATIC_DRAW) // *4 because  float32 is 4 bytes
	gl.EnableVertexAttribArray(uint32(POSITION_VB))
	gl.VertexAttribPointer(uint32(POSITION_VB), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	if len(m.quickmodel.UVs) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[TEXCOORD_VB])
		gl.BufferData(gl.ARRAY_BUFFER, len(m.quickmodel.UVs)*2*4, gl.Ptr(m.quickmodel.UVs), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(uint32(TEXCOORD_VB))
		gl.VertexAttribPointer(uint32(TEXCOORD_VB), 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	}

	if len(m.quickmodel.Normals) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo[NORMAL_VB])
		gl.BufferData(gl.ARRAY_BUFFER, len(m.quickmodel.Normals)*3*4, gl.Ptr(m.quickmodel.Normals), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(uint32(NORMAL_VB))
		gl.VertexAttribPointer(uint32(NORMAL_VB), 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	}

	gl.BindVertexArray(0)
}

func (m *Mesh) Draw() {
	if m.quickmodel != nil {
		m.DrawNew()
	} else {
		m.DrawOld()
	}
}

func (m *Mesh) DrawOld() {
	// defer util.TimeTrack(time.Now(), "mesh draw")
	gl.BindVertexArray(m.vao)

	// gl.DrawElements(gl.TRIANGLES, int32(len(m.model.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DrawElements(gl.TRIANGLES, int32(len(m.model.Indices))*3, gl.UNSIGNED_INT, gl.PtrOffset(0))
	// gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(len(m.model.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0), 0)

	gl.BindVertexArray(0)
}

func (m *Mesh) DrawNew() {
	// defer util.TimeTrack(time.Now(), "mesh quickdraw")
	gl.BindVertexArray(m.vao)

	// gl.DrawElements(gl.TRIANGLES, int32(len(m.model.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.quickmodel.Vertices))*3)
	// gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(len(m.model.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0), 0)

	gl.BindVertexArray(0)
}
