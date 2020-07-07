package errors

/*EmptyPostalToCordsError error to be raised when public.opendatasoft.com returns no coordinate data*/
type EmptyPostalToCordsError struct{}

func (m *EmptyPostalToCordsError) Error() string {
	return "The Query to the public.opendatasoft.com api returned no coordinate data."
}
