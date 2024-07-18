package filesearcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type searchfile struct {
	found   bool
	name    string
	dirName string
}

// channel that indicates the file was found
var fileFoundChan chan searchfile = make(chan searchfile)

func Init() {
	var wg sync.WaitGroup

	fmt.Println()
	fmt.Println("File Searcher")
	fmt.Println("-------------")

	if len(os.Args) < 2 {
		fmt.Println("Please enter target directory and search term.")
		return
	}

	dir := os.Args[1]

	searchTerm := os.Args[2]
	fmt.Println("Target Directory:", dir)
	fmt.Println("Search Term:", searchTerm)

	wg.Add(1)
	go func() {
		defer wg.Done() // remove counter when this concurrenly running thread of concurrently
		// running threads are done

		// keep this concurrently running until its done
		err := searchFiles(dir, searchTerm, &wg)

		if err != nil {
			log.Fatal("Error when searching files:", err)
		}
	}()

	// concurrently wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(fileFoundChan)
	}()

	result_file := <-fileFoundChan

	if result_file.found {
		fmt.Println("File was found.")
		fmt.Printf("File name: %s, directory name: %s", result_file.name, result_file.dirName)
	} else {
		fmt.Println("File was not found.")
	}
}

// return if file is found
func searchFiles(dir string, searchTerm string, wg *sync.WaitGroup) error {
	defer wg.Done()

	file, err := os.Open(dir)
	defer file.Close()

	if err != nil {
		fmt.Println("error while reading file:", err)
		return err
	}

	entries, err := file.ReadDir(-1)

	for _, entry := range entries {

		// skip if target is not a file but a directory
		if entry.IsDir() {
			wg.Add(1) // increment wait group before new goroutine

			// if its a directory we start a goroutine to search for the file there
			nestedPath := filepath.Join(dir, entry.Name())

			go searchFiles(nestedPath, searchTerm, wg)

			// skip the rest of the loop as we aren't matching folders
			continue
		}

		// found matching name
		fmt.Println()
		fmt.Println("SEARCHING IN DIRECTORY:", dir)
		fmt.Println("FILE NAME:", entry.Name())
		fmt.Println("SEARCH TERM:", searchTerm)
		fmt.Println()
		if entry.Name() == searchTerm {

			fileFoundChan <- searchfile{
				name:    entry.Name(),
				dirName: dir,
				found:   true,
			}

			// close channel after we found the target file
			close(fileFoundChan)
		}
	}
	if err != nil {
		fmt.Println("error while reading all directories")
		return err
	}

	return nil
}
