package filesearcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func Init() {
	fmt.Println()
	fmt.Println("File Searcher")
	fmt.Println("-------------")

	if len(os.Args) < 2 {
		fmt.Println("Please enter target directory and search term.")
		return
	}

	dir := os.Args[0]

	searchTerm := os.Args[1]
	fmt.Println("Target Directory:", dir)
	fmt.Println("Search Term:", searchTerm)

	found, err := searchFiles(dir, searchTerm)

	if err != nil {
		log.Fatal("Error when searching files:", err)
	}

	fmt.Println("File was found:", found)
}

// return if file is found
func searchFiles(dir string, searchTerm string) (bool, error) {

	file, err := os.Open(dir)

	if err != nil {
		fmt.Println("error while reading file:", err)
		return false, err
	}

	defer file.Close()

	entries, err := file.ReadDir(-1)

	if err != nil {
		fmt.Println("error while reading all directories")
		return false, err
	}

	for _, entry := range entries {
		var wg *sync.WaitGroup
		// skip if target is not a file but a directory
		if entry.IsDir() {
			// if its a directory we start a goroutine to search for the file there
			nestedPath := filepath.Join(dir, entry.Name())

			wg.Add(1)
			go searchFiles(nestedPath, searchTerm)
			wg.Wait()
			continue
		}

		fmt.Print("entry.Name:", entry.Name())
		fmt.Println(" compared to search term:", searchTerm)
		if entry.Name() == searchTerm {
			wg.Done()
			return true, nil
		}

	}

	return false, nil
}
