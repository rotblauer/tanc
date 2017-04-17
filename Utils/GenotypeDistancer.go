package Utils

//VectorDistancer A Distancer implemented for vectors of float64
type GenotypeDistancer struct {
	Matrix [][]float64
}

// Len returns the length of the vector.
// This function is needed so that VectorDistancer implements the
// Distancer interface.
func (vd GenotypeDistancer) Len() int { return len(vd.Matrix) }

// Distance returns the euclidean distance between vd[i] and vd[j].
// This function is needed so that VectorDistancer implements the
// Distancer interface.
func (vd GenotypeDistancer) Distance(i, j int) float64 {
	vi := vd.Matrix[i]
	vj := vd.Matrix[j]
	dist := 0.0
	for k, vik := range vi {
		vjk := vj[k]
		dist += (vik - vjk) * (vik - vjk)
	}
	return dist
}
