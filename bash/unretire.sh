#moves all related content for a particular list to the retired folder.

#note the underscore forces exact names only
mv $VXDAY_RETIRED_DIR/$1_*.vxday $VXDAY_ACTIVE_DIR
