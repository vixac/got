package bullet_engine

//VX:TODO
/*
type ListStore struct {
	Mesh bullet_stl.Mesh
}

type ListId struct {
	Id string
}
type ListStoreInterface interface {
	AddToList(list ListId, ids []int32) error
	All() (map[ListId][]int32, error)
	ListNames() ([]string, error)
	UpdateName(list ListId, name string) error
	RemoveFromList(ids []int32) error
}

func NewList(client bullet_interface.TrackClientInterface, meshName string, bucketId int32, subjectSeparor string, forwardSeparator string, backwardSeparator string) (ListStoreInterface, error) {

	mesh, err := bullet_stl.NewBulletMesh(client, bucketId, meshName, forwardSeparator, backwardSeparator)
	if err != nil {
		return nil, err
	}

	return &ListStore{
		Mesh: mesh,
	}, nil
}

func (l *ListStore) AddToList(ids []int32) error {
	var pairs []bullet_stl.ManyToManyPair
	return errors.New("not impl")
}
func (l *ListStore) All() ([]int32, error) {
	return []int32{}, errors.New("not impl")
}
func (l *ListStore) ListName() (string, error) {
	return "", errors.New("not impl")
}
func (l *ListStore) UpdateName(name string) error {
	return errors.New("not impl")
}
func (l *ListStore) RemoveFromList(ids []int32) error {
	return errors.New("not impl")
}
*/
