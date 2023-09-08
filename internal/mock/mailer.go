package mock

type MockMailer struct{}

func (m MockMailer) Send(recipient, templateFile string, data interface{}) error {
	return nil
}
