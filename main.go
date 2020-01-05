package main

func main() {
	// menu()
	s1 := []movie{getMovie("tt0314331"), getMovie("tt0095016")}
	s2 := []movie{getMovie("tt3748528"), getMovie("tt2283362")}
	createHTML(s1, s2)
}
