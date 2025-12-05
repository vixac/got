rm got-bolt

echo "----- go run it..."
go run . it Top Item.

echo "----- go run alias..."
go run . alias -a top -g 0a0
echo "----- go run note top level2...."
go run . note top level2 
echo "----- go run alias l2..."
go run . alias -a l2 -g 0a1
echo "----- go run note l2 level3"
go run . note l2 level3
echo "----- go run note l2 level4"
go run . note l2 level4
echo "----- go run note l2 level5"
go run . note l2 level5
echo "----- go it Top Leaf."
go run . it Top leaf.
echo "----- go run jobs"
go run . jobs
