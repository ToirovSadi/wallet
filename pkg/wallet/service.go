package wallet

func (s *Service) Import(dir string) (err error) {
	err = s.importAccounts(dir + "\\accounts.dump")
	if err != nil {
		return err
	}

	err = s.importPayments(dir + "\\payments.dump")
	if err != nil {
		return err
	}
	err = s.importFavorites(dir + "\\favorites.dump")
	if err != nil {
		return err
	}
	return err
}

func (s *Service) Export(dir string) (err error) {
	err = exportAccounts(s.accounts, dir+"\\accounts.dump")
	if err != nil {
		return err
	}
	err = exportPayments(s.payments, dir+"\\payments.dump")
	if err != nil {
		return err
	}
	err = exportFavorites(s.favorites, dir+"\\favorites.dump")
	return err
}
