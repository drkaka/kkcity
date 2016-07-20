package kkcity

import "github.com/stretchr/testify/suite"

type languageHandleSuite struct {
	suite.Suite
}

func (suite *languageHandleSuite) TestGetLanguage() {
	var lang0, lang1 string
	var err error

	lang0, err = getLanguage(0)
	suite.NoError(err, "Should be able to get language.")
	suite.Equal(testLangs[0], lang0, "Language at index 0 is wrong.")

	lang1, err = getLanguage(1)
	suite.NoError(err, "Should be able to get language.")
	suite.Equal(testLangs[1], lang1, "Language at index 0 is wrong.")

	_, err = getLanguage(2)
	suite.Equal(ErrLanguageIndex, err, "Language will be out of range.")
}

func (suite *languageHandleSuite) TestGetAllLanguage() {
	all := getAll()
	suite.EqualValues(testLangs, all, "Get all languages is wrong.")
}
