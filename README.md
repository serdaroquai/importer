# importer
golang csv to csv transformer for n-ary tree representations. Input csv is in naive unsorted parent-child representation. Output csv is a closure table format which includes depth information of a child node to each of its parents

## input format (descendant,ancestor)
C,A  
D,C  
E,A  
F,C  
G,F  
H,E  

## output format (ancestor, descendant, depth)
A,A,0  
C,C,0  
A,C,1  
E,E,0  
A,E,1  
D,D,0  
C,D,1  
A,D,2  
F,F,0  
C,F,1  
A,F,2  
H,H,0  
E,H,1  
A,H,2  
G,G,0  
F,G,1  
C,G,2  
A,G,3  

# Usage

`go get github.com/serdaroqui/importer`

`go run import.go <inputPath> <outputPath>`
