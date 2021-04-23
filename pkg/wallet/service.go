package wallet

func (s *Service) Import(dir string) (err error) {
	err = s.importAccounts(dir + "/accounts.dump")
	if err != nil {
		return err
	}
	err = s.importPayments(dir + "/payments.dump")
	if err != nil {
		return err
	}
	err = s.importFavorites(dir + "/favorites.dump")
	if err != nil {
		return err
	}
	return err
}

func (s *Service) Export(dir string) (err error) {
	err = s.exportAccounts(dir + "\\accounts.dump")
	if err != nil {
		return err
	}
	err = s.exportPayments(dir + "\\payments.dump")
	if err != nil {
		return err
	}
	err = s.exportFavorites(dir + "\\favorites.dump")
	return err
}
