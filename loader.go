package main

import (
	"bufio"
	"fmt"
	"os"

	g "github.com/lambertwang/tracery/geometry"
	sc "github.com/lambertwang/tracery/scene"
)

type face struct {
	vec0, vec1, vec2    int
	tex0, tex1, tex2    int
	norm0, norm1, norm2 int
}

type model struct {
	vecs, uvs, norms []g.Vector
	faces            []face
}

func loadObjModel(fileName string) (outModel model) {
	objFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer objFile.Close()

	scanner := bufio.NewScanner(objFile)
	scanner.Split(bufio.ScanLines)

	outModel = model{}

	for scanner.Scan() {
		line := scanner.Text()
		var prefix string

		fmt.Sscanf(line, "%s", &prefix)

		switch prefix {
		case "v":
			v := g.Vector{}
			fmt.Sscanf(line, "%s %f %f %f\n", &prefix,
				&v.X, &v.Y, &v.Z)
			outModel.vecs = append(outModel.vecs, v)
			break
		case "vt":
			v := g.Vector{}
			fmt.Sscanf(line, "%s %f %f %f\n", &prefix,
				&v.X, &v.Y)
			outModel.uvs = append(outModel.uvs, v)
			break
		case "vn":
			v := g.Vector{}
			fmt.Sscanf(line, "%s %f %f %f\n", &prefix,
				&v.X, &v.Y, &v.Z)
			outModel.norms = append(outModel.norms, v)
			break
		case "f":
			vecI := make([]int, 3)
			uvI := make([]int, 3)
			normI := make([]int, 3)

			count, _ := fmt.Sscanf(line, "%s %d/%d/%d %d/%d/%d %d/%d/%d", &prefix,
				&vecI[0], &uvI[0], &normI[0],
				&vecI[1], &uvI[1], &normI[1],
				&vecI[2], &uvI[2], &normI[2],
			)
			if count == 10 {
				outModel.faces = append(outModel.faces, face{
					vecI[0] - 1, vecI[1] - 1, vecI[2] - 1,
					uvI[0] - 1, uvI[1] - 1, uvI[2] - 1,
					normI[0] - 1, normI[1] - 1, normI[2] - 1,
				})
				break
			}

			count, _ = fmt.Sscanf(line, "%s %d//%d %d//%d %d//%d", &prefix,
				&vecI[0], &normI[0],
				&vecI[1], &normI[1],
				&vecI[2], &normI[2],
			)
			if count == 7 {
				outModel.faces = append(outModel.faces, face{
					vecI[0] - 1, vecI[1] - 1, vecI[2] - 1,
					-1, -1, -1,
					normI[0] - 1, normI[1] - 1, normI[2] - 1,
				})
				break
			}
			count, _ = fmt.Sscanf(line, "%s %d/%d %d/%d %d/%d", &prefix,
				&vecI[0], &uvI[0],
				&vecI[1], &uvI[1],
				&vecI[2], &uvI[2],
			)
			if count == 7 {
				outModel.faces = append(outModel.faces, face{
					vecI[0] - 1, vecI[1] - 1, vecI[2] - 1,
					uvI[0] - 1, uvI[1] - 1, uvI[2] - 1,
					-1, -1, -1,
				})
				break
			}

			count, _ = fmt.Sscanf(line, "%s %d %d %d", &prefix,
				&vecI[0], &vecI[1], &vecI[2],
			)
			if count == 4 {
				outModel.faces = append(outModel.faces, face{
					vecI[0] - 1,
					vecI[1] - 1,
					vecI[2] - 1,
					-1, -1, -1, -1, -1, -1,
				})
			}
			panic("Unable to read face data")
			break
		case "#":
		default:
			break
		}
	}

	return
}

func (m model) toShapes(mat sc.Material, t g.TMatrix3d) []sc.Triangle {
	var outShapes []sc.Triangle

	for _, f := range m.faces {
		newTri := sc.CreateTriangle(
			mat,
			t.Transform(m.vecs[f.vec0]),
			t.Transform(m.vecs[f.vec1]),
			t.Transform(m.vecs[f.vec2]),
		)
		if f.norm0 >= 0 && f.norm1 >= 0 && f.norm2 >= 0 {
			newTri.VertexNormals[0] = t.Transform(m.norms[f.norm0]).Norm()
			newTri.VertexNormals[1] = t.Transform(m.norms[f.norm1]).Norm()
			newTri.VertexNormals[2] = t.Transform(m.norms[f.norm2]).Norm()
		}

		outShapes = append(outShapes, newTri)
	}
	fmt.Printf("Loaded %d tris\n", len(outShapes))

	return outShapes
}
