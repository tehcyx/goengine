package obj

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tehcyx/goengine/util"
)

type QuickObjModel struct {
	Vertices []mgl32.Vec3
	UVs      []mgl32.Vec2
	Normals  []mgl32.Vec3
}

func LoadObj(path string) *QuickObjModel {
	defer util.TimeTrack(time.Now(), "loadObj")
	tmp_model := new(QuickObjModel)

	var vertexIndices, uvIndices, normalIndices []int

	fileHandle, _ := os.Open(path)
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
			tmp_model.Vertices = append(tmp_model.Vertices, vec)
		//INDICES
		case "f":
			tokens := strings.Split(line, " ")
			tokens = delete_empty(tokens)

			vertex1 := strings.Split(tokens[1], "/")
			vertex2 := strings.Split(tokens[2], "/")
			vertex3 := strings.Split(tokens[3], "/")

			vertexIndices = append(vertexIndices, parseIndexValue(vertex1[0]))
			vertexIndices = append(vertexIndices, parseIndexValue(vertex1[1]))
			vertexIndices = append(vertexIndices, parseIndexValue(vertex1[2]))

			uvIndices = append(uvIndices, parseIndexValue(vertex2[0]))
			uvIndices = append(uvIndices, parseIndexValue(vertex2[1]))
			uvIndices = append(uvIndices, parseIndexValue(vertex2[2]))

			normalIndices = append(normalIndices, parseIndexValue(vertex3[0]))
			normalIndices = append(normalIndices, parseIndexValue(vertex3[1]))
			normalIndices = append(normalIndices, parseIndexValue(vertex3[2]))
		// NORMALS.
		case "vn":
			vec := parseVec3(line)
			tmp_model.Normals = append(tmp_model.Normals, vec)
		// TEXTURE VERTICES.
		case "vt":
			vec := parseVec2(line)
			tmp_model.UVs = append(tmp_model.UVs, vec)
		}
	}

	result := new(QuickObjModel)
	for i := 0; i < len(vertexIndices); i++ {
		vertexIndex := vertexIndices[i]
		if vertexIndex-1 < len(tmp_model.Vertices) && vertexIndex-1 >= 0 {
			// fmt.Printf("len indices: %d\n", len(vertexIndices))
			// fmt.Printf("len uvs: %d\n", len(tmp_model.Vertices))
			// fmt.Printf("[i]: %d\n", vertexIndices[i])
			// fmt.Printf("i: %d\n", i)
			result.Vertices = append(result.Vertices, tmp_model.Vertices[vertexIndex-1])
		}

		if len(uvIndices) > 0 {
			uvIndex := uvIndices[i]
			if uvIndex-1 < len(tmp_model.UVs) && uvIndex-1 >= 0 {
				result.UVs = append(result.UVs, tmp_model.UVs[uvIndex-1])
			}
		}

		if len(normalIndices) > 0 {
			normalIndex := normalIndices[i]
			if normalIndex-1 < len(tmp_model.Normals) && normalIndex-1 >= 0 {
				result.Normals = append(result.Normals, tmp_model.Normals[normalIndex-1])
			}
		}
	}
	return result
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
