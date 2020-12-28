package dbservice

// Select======================
var getMissingDataInput = "SELECT * FROM daily_raw_data t1 WHERE NOT EXISTS (SELECT * FROM data_input t2 WHERE t1.symbol = t2.symbol and t1.dt = t2.dt)"

func (s *DBService) GetMissingDataInput(data *[]RawDataEntity) error {
	_, err := s.client.DB.Select(data, getMissingDataInput)
	if err != nil {
		return err
	}
	return nil
}

// insert ==
func (s *DBService) InsertDataInput(inputData []*DataInput) int {
	count := 0
	for _, v := range inputData {
		err := s.client.DB.Insert(v)
		if err == nil {
			count++
		}
	}
	return count
}
