package main

import (
	"time"
	"encoding/csv"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/serdaroquai/importer/ntree"
)

const (
	self = iota
	parent
)

func main() {

	defer elapsed()() //time the whole thing

	args := os.Args[1:] // input csv args[0], output csv args[1]

	file, err := os.Open(args[0])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	// read records of connectivity
	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	// create nodes
	nodes := createNodes(rows)

	// find out root nodes
	var rootNodes []*ntree.Tree
	for _, node := range nodes {
		if node.IsRoot() {
			rootNodes = append(rootNodes, ntree.NewTree(node))
		}
	}

	// create csvWriter
	writerCh, outputFile, writerDoneCh := createCsvWriter(args[1])
	defer func() {
		outputFile.Close()
		log.Printf("Closed output file %v", args[1])
	}()

	// find out number of available CPUs
	workerCount :=  runtime.NumCPU()
	log.Printf("Available Cpu count %v\n", workerCount)

	// create workers
	workCh := make(chan *ntree.Tree)
	doneCh := make(chan bool)
	for i := 0; i < workerCount; i++ {
		createWorker(i, workCh, doneCh, writerCh)
	}

	// distribute work
	for _, tree := range rootNodes {
		workCh <- tree
	}
	// signal all workers done
	for i := 0; i < workerCount; i++ {
		doneCh <- true
	}

	// signal writer done
	writerDoneCh <- true
}

func elapsed() func() {
    start := time.Now()
    return func() {
        log.Printf("It took %v\n", time.Since(start))
    }
}

func createWorker(id int, submitWorkCh chan *ntree.Tree, doneCh chan bool, csvWriterCh chan [][]string) {
	go func() {
		workCount := 0

		log.Printf("worker %v started", id)
		for done := false; !done; {
			select {
			case tree := <-submitWorkCh:
				var result [][]string
				
				tree.BreadthFirst(func(node *ntree.Node) {
					//first apply self as depth 0
					result = append(result, []string{node.Id, node.Id, "0"})

					depth := 1
					for parent := node.Parent; parent != nil; {
						// then apply a line for each parent
						result = append(result, []string{node.Id, parent.Id, strconv.Itoa(depth)})
						// update depth and new parent
						depth++
						parent = parent.Parent
					}
				})

				// send results
				workCount++
				csvWriterCh <- result

			case done = <-doneCh:
				// end worker
				log.Printf("worker %v ended with a workcount of %v", id, workCount)
			}
		}
	}()
}

func createCsvWriter(filename string) (inputCh chan [][]string, file *os.File, doneCh chan bool) {
	inputCh = make(chan [][]string) // used to receive chunks of input
	doneCh = make(chan bool)        // used to make sure writer is finished

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(file)

	go func() {
		for done:=false; !done; {
			select {
			case input := <-inputCh:
				//todo write and flush input
				if err := w.WriteAll(input); err != nil {
					log.Fatalln("Error writing record to csv:", err)
				}
			case done = <-doneCh:
				// just receive this to let main func know writing is done
			}
		
		}
	}()
	return
}

func createNodes(rows [][]string) map[string]*ntree.Node {
	nodes := make(map[string]*ntree.Node)
	for _, row := range rows {

		selfNode, present := nodes[row[self]] // check if node exists already

		if !present { // create node if not present
			selfNode = ntree.NewNode(row[self])
			nodes[row[self]] = selfNode
		}

		parentNode, present := nodes[row[parent]] // check if parent node exists already

		if !present { // create node if not present
			parentNode = ntree.NewNode(row[parent])
			nodes[row[parent]] = parentNode
		}

		if row[parent] != "" { // add to parent if it has a parent
			nodes[row[parent]].AddChild(selfNode)
		}
	}
	return nodes
}
