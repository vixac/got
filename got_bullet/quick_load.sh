rm got-bolt

echo "----- go run it..."
go run . it Top Item.

echo "----- go run alias..."
go run . alias top 0a0
echo "----- go run note top level2...."
go run . note top level2 
echo "----- go run alias l2..."
go run . alias Level2 0a1
echo "----- go run note l2 1st Note "
go run . note Level2 This is a note about The Level 2 task
echo "----- go run note l2 level4"
go run . note Level2 This is a second note
echo "----- go run note l2 This is another note for Level 2"
go run . under Level2 level3 task
echo "----- go it Top Leaf."
go run . it Top leaf.
echo "----- go run jobs"
go run . jobs
