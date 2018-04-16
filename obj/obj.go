package obj

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type ObjectIndex struct {
	VertexIndex int
	UVIndex     int
	NormalIndex int
}

type IndexedModel struct {
	Positions []mgl32.Vec3
	TexCoords []mgl32.Vec2
	Normals   []mgl32.Vec3
	Indices   []int
}

type ObjModel struct {
	Indices  []ObjectIndex
	Vertices []mgl32.Vec3
	UVs      []mgl32.Vec2
	Normals  []mgl32.Vec3
}

func NewObjModelFromFile(filePath string) *ObjModel {
	o := new(ObjModel)

	fileHandle, _ := os.Open(filePath)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		var lineType string
		if len(line) < 2 {
			continue
		}
		lineType = strings.Fields(line)[0]

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			vec := parseVec3(line)
			o.Vertices = append(o.Vertices, vec)
		//INDICES
		case "f":
			o.CreateFace(line)
		// NORMALS.
		case "vn":
			vec := parseVec3(line)
			o.Normals = append(o.Normals, vec)
		// TEXTURE VERTICES.
		case "vt":
			vec := parseVec2(line)
			o.UVs = append(o.UVs, vec)
		}
	}

	return o
}

func (o *ObjModel) ToIndexedModel() *IndexedModel {
	result := new(IndexedModel)
	normalModel := new(IndexedModel)

	var indexLookup []ObjectIndex

	copy(indexLookup, o.Indices)

	// sort indexes
	sort.Slice(indexLookup[:], func(i, j int) bool {
		return indexLookup[i].VertexIndex < indexLookup[j].VertexIndex
	})

	var normalModelIndexMap map[ObjectIndex]int
	normalModelIndexMap = make(map[ObjectIndex]int)
	var indexMap map[int]int
	indexMap = make(map[int]int)

	for i := 0; i < len(o.Indices); i++ {
		currentIndex := &o.Indices[i]
		var currentPosition mgl32.Vec3
		var currentTexCoord mgl32.Vec2
		var currentNormal mgl32.Vec3

		if o.HasUVs() {
			currentTexCoord = o.UVs[currentIndex.UVIndex]
		} else {
			currentTexCoord = mgl32.Vec2{0.0, 0.0}
		}

		if o.HasNormals() {
			currentNormal = o.Normals[currentIndex.NormalIndex]
		} else {
			currentNormal = mgl32.Vec3{0.0, 0.0, 0.0}
		}

		var normalModelIndex int
		var resultModelIndex int

		//Create model to properly generate Normals on
		if val, ok := normalModelIndexMap[*currentIndex]; ok {
			normalModelIndex = int(val)

			normalModelIndexMap[*currentIndex] = normalModelIndex
			normalModel.Positions = append(normalModel.Positions, currentPosition)
			normalModel.TexCoords = append(normalModel.TexCoords, currentTexCoord)
			normalModel.Normals = append(normalModel.Normals, currentNormal)
		} else {
			normalModelIndex = len(normalModel.Positions)
		}

		//Create model which properly separates texture coordinates
		previousVertexLocation := o.findLastVertexIndex(indexLookup, *currentIndex, result)

		if previousVertexLocation == -1 {
			resultModelIndex = len(result.Positions)
			result.Positions = append(result.Positions, currentPosition)
			result.TexCoords = append(result.TexCoords, currentTexCoord)
			result.Normals = append(result.Normals, currentNormal)
		} else {
			resultModelIndex = previousVertexLocation
		}

		normalModel.Indices = append(normalModel.Indices, normalModelIndex)
		result.Indices = append(result.Indices, resultModelIndex)
		indexMap[resultModelIndex] = normalModelIndex
	}

	if !o.HasNormals() {
		normalModel.CalcNormals()

		for j := 0; j < len(result.Positions); j++ {
			result.Normals[j] = normalModel.Normals[indexMap[j]]
		}
	}

	return result
}

func (o *ObjModel) findLastVertexIndex(indexLookup []ObjectIndex, currentIndex ObjectIndex, result *IndexedModel) int {
	start := 0
	end := len(indexLookup)
	current := (end-start)/2 + start
	previous := start

	for current != previous {
		testIndex := indexLookup[current]

		if testIndex.VertexIndex == currentIndex.VertexIndex {
			countStart := current

			for i := 0; i < current; i++ {
				possibleIndex := indexLookup[current-i]

				if possibleIndex == currentIndex {
					continue
				}

				if possibleIndex.VertexIndex != currentIndex.VertexIndex {
					break
				}

				countStart++
			}

			for i := countStart; i < len(indexLookup)-countStart; i++ {
				possibleIndex := indexLookup[current+i]

				if possibleIndex == currentIndex {
					continue
				}

				if possibleIndex.VertexIndex != currentIndex.VertexIndex {
					break
				} else if (!o.HasUVs() || possibleIndex.UVIndex == currentIndex.UVIndex) && (!o.HasNormals() || possibleIndex.NormalIndex == currentIndex.NormalIndex) {
					currentPosition := o.Vertices[currentIndex.VertexIndex]
					var currentTexCoord mgl32.Vec2
					var currentNormal mgl32.Vec3

					if o.HasUVs() {
						currentTexCoord = o.UVs[currentIndex.UVIndex]
					} else {
						currentTexCoord = mgl32.Vec2{0.0, 0.0}
					}

					if o.HasNormals() {
						currentNormal = o.Normals[currentIndex.NormalIndex]
					} else {
						currentNormal = mgl32.Vec3{0.0, 0.0, 0.0}
					}

					for j := 0; j < len(result.Positions); j++ {
						if currentPosition == result.Positions[j] &&
							(!o.HasUVs() || currentTexCoord == result.TexCoords[j]) &&
							(!o.HasNormals() || currentNormal == result.Normals[j]) {
							return j
						}
					}
				}
			}
			return -1
		} else {
			if testIndex.VertexIndex < currentIndex.VertexIndex {
				start = current
			} else {
				end = current
			}
		}
		previous = current
		current = (end-start)/2 + start
	}

	return -1
}

func (o *ObjModel) HasUVs() bool {
	if len(o.UVs) > 0 {
		return true
	}
	return false
}

func (o *ObjModel) HasNormals() bool {
	if len(o.Normals) > 0 {
		return true
	}
	return false
}

func (o *ObjModel) CreateFace(line string) {
	tokens := strings.Split(line, " ")

	o.Indices = append(o.Indices, o.parseIndex(tokens[1]))
	o.Indices = append(o.Indices, o.parseIndex(tokens[2]))
	o.Indices = append(o.Indices, o.parseIndex(tokens[3]))

	if len(tokens) > 4 {
		o.Indices = append(o.Indices, o.parseIndex(tokens[1]))
		o.Indices = append(o.Indices, o.parseIndex(tokens[3]))
		o.Indices = append(o.Indices, o.parseIndex(tokens[4]))
	}
}

func parseVec2(line string) mgl32.Vec2 {
	vertexIndexStart := 3

	for vertexIndexStart < len(line) {
		if string(line[vertexIndexStart]) != " " {
			break
		}
		vertexIndexStart++
	}

	vertexIndexEnd := findNextChar(vertexIndexStart, line, " ")

	x := parseFloatValue(line, vertexIndexStart, vertexIndexEnd)

	vertexIndexStart = vertexIndexEnd + 1
	vertexIndexEnd = findNextChar(vertexIndexStart, line, " ")

	y := parseFloatValue(line, vertexIndexStart, vertexIndexEnd+1)

	return mgl32.Vec2{x, y}
}

func parseVec3(line string) mgl32.Vec3 {
	vertexIndexStart := 2

	for vertexIndexStart < len(line) {
		if string(line[vertexIndexStart]) != " " {
			break
		}
		vertexIndexStart++
	}

	vertexIndexEnd := findNextChar(vertexIndexStart, line, " ")

	x := parseFloatValue(line, vertexIndexStart, vertexIndexEnd)

	vertexIndexStart = vertexIndexEnd + 1
	vertexIndexEnd = findNextChar(vertexIndexStart, line, " ")

	y := parseFloatValue(line, vertexIndexStart, vertexIndexEnd)

	vertexIndexStart = vertexIndexEnd + 1
	vertexIndexEnd = findNextChar(vertexIndexStart, line, " ")

	z := parseFloatValue(line, vertexIndexStart, vertexIndexEnd+1)

	return mgl32.Vec3{x, y, z}
}

func (o *ObjModel) parseIndex(token string) ObjectIndex {
	var result ObjectIndex

	vertexIndexStart := 0
	vertexIndexEnd := findNextChar(vertexIndexStart, token, "/")

	result.VertexIndex = parseIndexValue(token, vertexIndexStart, vertexIndexEnd)
	result.UVIndex = 0
	result.NormalIndex = 0

	if vertexIndexEnd >= len(token) {
		return result
	}

	vertexIndexStart = vertexIndexEnd + 1
	vertexIndexEnd = findNextChar(vertexIndexStart, token, "/")

	result.UVIndex = parseIndexValue(token, vertexIndexStart, vertexIndexEnd)

	if vertexIndexEnd >= len(token) {
		return result
	}

	vertexIndexStart = vertexIndexEnd + 1
	vertexIndexEnd = findNextChar(vertexIndexStart, token, "/")

	result.NormalIndex = parseIndexValue(token, vertexIndexStart, vertexIndexEnd)

	return result
}

func findNextChar(start int, token, find string) int {
	result := start
	for result < len(token)-1 {
		result++

		if string(token[result]) == find {
			break
		}
	}
	return result
}

func parseIndexValue(token string, vertexIndexStart, vertexIndexEnd int) int {
	res, _ := strconv.Atoi(token[vertexIndexStart:vertexIndexEnd])
	return res
}

func parseFloatValue(token string, vertexIndexStart, vertexIndexEnd int) float32 {
	res, _ := strconv.ParseFloat(token[vertexIndexStart:vertexIndexEnd], 32)
	return float32(res)
}

func (im *IndexedModel) CalcNormals() {
	for i := 0; i < len(im.Indices); i += 3 {
		i0 := im.Indices[i]
		i1 := im.Indices[i+1]
		i2 := im.Indices[i+2]

		v1 := im.Positions[i1].Sub(im.Positions[i0])
		v2 := im.Positions[i2].Sub(im.Positions[i0])
		v3 := (v1.Cross(v2)).Normalize()

		im.Normals[i0].Add(v3)
		im.Normals[i1].Add(v3)
		im.Normals[i2].Add(v3)
	}

	for i := 0; i < len(im.Positions); i++ {
		im.Normals[i] = im.Normals[i].Normalize()
	}
}
