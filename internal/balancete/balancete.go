package balancete

import (
	chartofaccount "Clarion/internal/chartOfAccount"
	"Clarion/internal/models"
	"Clarion/internal/movement"
	"fmt"
	"sort"
	"strconv"
)

/* Função criada por Ricardo Silva Ferreira
   Inicio da criação 19/03/2025 14:00
   Data Final da criação : 19/03/2025 15:00
*/

func VerifyExistAccount(codAccount string, matriz []models.Balancete) (bool, int) {

	for a := 0; a < len(matriz); a++ {

		if codAccount == matriz[a].CodAccount {
			return true, a
		}
	}

	return false, -1
}

/*
	  	Essa função verifica todas as contas e faz a soma até chegar ao pai raiz
		exemplo se a conta é 511 e o valor de débito é 20 ela leva esse valor para o débito da conta 51 que é o pai
*/

/*
Função criada por Ricardo Silva Ferreira

	Inicio da criação 21/03/2025 12:00
	Data Final da criação : 21/03/2025 12:40
*/
func SumValuesFathers(codAccount string, _balancetes []models.Balancete) {
	totolreg := len(codAccount)

	debit := 0.0
	credit := 0.0

	for totolreg >= 2 {

		codAccountFather := codAccount[0:totolreg]

		for t := len(_balancetes) - 1; t >= 0; t-- {

			//Aqui guarda o valor da primeira conta com valor ou seja a conta que tem valor
			if _balancetes[t].CodAccount == codAccount && debit == 0 && credit == 0 {
				debit = _balancetes[t].DebitValue
				credit = _balancetes[t].CreditValue
				fmt.Println("conta com valor", codAccountFather, debit, credit)

			} else if _balancetes[t].CodAccount == codAccountFather {

				fmt.Println("codAccountFather pai", codAccountFather, debit, credit)

				_balancetes[t].DebitValue = _balancetes[t].DebitValue + debit
				_balancetes[t].CreditValue = _balancetes[t].CreditValue + credit
				// debit = 0.0
				// credit = 0.0
			}

		}

		totolreg--
	}
	fmt.Println("Acabou a conta ", codAccount)

}

/* Função criada por Ricardo Silva Ferreira
   Inicio da criação 21/03/2025 12:40
   Data Final da criação : 21/03/2025 12:45
*/

func ReturnBalanceteLine(codAccount string, debitValue float64, creditValue float64, description string, fatherCod string, class string, sum bool) models.Balancete {

	var mvi models.Balancete

	mvi.CodAccount = codAccount
	mvi.DebitValue = debitValue
	mvi.CreditValue = creditValue
	mvi.Description = description
	mvi.FatherCod = fatherCod
	mvi.Class = class
	mvi.Sum = sum
	return mvi
}

/* Função criada por Ricardo Silva Ferreira
   Inicio da criação 19/03/2025 12:00
   Data Final da criação : 21/03/2025 13:09
*/

func GenerateBalanceteReport(initialYear int, initialMonth int, endYear int, endMonth int) []models.Balancete {

	movement, ert := movement.SearchRecordsBetweenMonths(initialYear, initialMonth, endYear, endMonth)

	var _balancetes []models.Balancete

	if ert == nil {

		for i := 0; i < len(movement); i++ {

			for j := 0; j < len(movement[i].Movements); j++ {

				var _accounts []string

				// var _balanceteTemp []models.Balancete

				// fmt.Println("Começo do balancete " + strconv.Itoa(i))

				for k := 0; k < len(movement[i].Movements[j].MovementsItens); k++ {

					movimentItensLine := movement[i].Movements[j].MovementsItens[k]

					if movimentItensLine.Active { //Apenas movimentos ativos

						existBalanceteRecord := false
						existAccount := false

						//O For abaixo verifica se já existe conta cadastrada para esse movimento
						for m := 0; m < len(_accounts); m++ {
							if _accounts[m] == movimentItensLine.CodAccount {
								existAccount = true
								break
							}
						}

						//Se não existir conta ele cadastra na matriz de contas
						if !existAccount {
							_accounts = append(_accounts, movimentItensLine.CodAccount)
						}

						for n := 0; n < len(_balancetes); n++ {

							//SE EXISTIR APENAS ATUALIZA OS VALORES COM A SOMA DE DÉBITO OU CRÉDITO
							if _balancetes[n].CodAccount == movimentItensLine.CodAccount {

								totalcredit := 0.00
								totalDebit := 0.0
								existBalanceteRecord = true

								//DADOS DE DÉBITO
								debit := _balancetes[n].DebitValue
								debitDB, _ := strconv.ParseFloat(movimentItensLine.DebitValue.String(), 64)
								totalDebit = debit + debitDB
								_balancetes[n].DebitValue = totalDebit
								_balancetes[n].Sum = true

								//DADOS DE CRÉDITO
								credit := _balancetes[n].CreditValue
								creditDB, _ := strconv.ParseFloat(movimentItensLine.CreditValue.String(), 64)
								totalcredit = credit + creditDB
								_balancetes[n].CreditValue = totalcredit

								break
							}
						}

						if !existBalanceteRecord {

							debit, _ := strconv.ParseFloat(movimentItensLine.DebitValue.String(), 64)
							credit, _ := strconv.ParseFloat(movimentItensLine.CreditValue.String(), 64)

							_balancetes = append(_balancetes, ReturnBalanceteLine(
								movimentItensLine.CodAccount,
								debit,
								credit,
								"",
								movimentItensLine.CodAccount[0:2],
								movimentItensLine.CodAccount[0:1],
								true,
							))

						}
					}
				}

				for o := 0; o < len(_accounts); o++ {

					/*
					  Pega a conta e vai subtraindo digitos até chegar na conta razão que possui apenas dois digitos
					  Ex: começa com a conta 51234 e vai subtrainto 5123 depos 512 e para no 51 e assim retorna uma Matriz
					  de 3 contas a 51 é a conta pai
					*/
					accountListBreak := chartofaccount.BreakAccounts(_accounts[o])

					/*
						Pega o resultado da variável accountListBreak e vai ao banco de dados e retorna todos os dados de cada conta
						o mais importante aqui é pegar as descrições
					*/
					accountListReady, err := chartofaccount.SearchAccounts(accountListBreak)

					if err == nil {

						for p := 0; p < len(accountListReady); p++ {

							exist, indice := VerifyExistAccount(accountListReady[p].CodAccount, _balancetes)
							if !exist {

								_balancetes = append(_balancetes, ReturnBalanceteLine(
									accountListReady[p].CodAccount,
									0,
									0,
									accountListReady[p].Description,
									accountListReady[p].CodAccount[0:2],
									accountListReady[p].CodAccount[0:1],
									false,
								))

								// fmt.Println("CodsAccount", mvi.CodAccount+" - "+mvi.Description+" - "+fmt.Sprintf("%.2f", mvi.DebitValue)+" - "+fmt.Sprintf("%.2f", mvi.CreditValue))

							} else {
								_balancetes[indice].Description = accountListReady[p].Description

								// fmt.Println("CodsAccount", _balanceteTemp[indice].CodAccount+" - "+_balanceteTemp[indice].Description+" - "+fmt.Sprintf("%.2f", _balanceteTemp[indice].DebitValue)+" - "+fmt.Sprintf("%.2f", _balanceteTemp[indice].CreditValue))

							}

						}
					}

					// SumValuesFathers(_accounts[o], _balancetes)

					// for t := 0; t < len(_balancetes); t++ {
					// 	if _balancetes[t].CodAccount[0:2] == "81" {
					// 		fmt.Println(_balancetes[t].CodAccount, " - ", _balancetes[t].Description, _balancetes[t].DebitValue, _balancetes[t].CreditValue, _balancetes[t].Sum)
					// 	}
					// }

				}

			}
		}

		sort.Slice(_balancetes, func(i, j int) bool {
			return _balancetes[i].CodAccount < _balancetes[j].CodAccount
		})

		for t := 0; t < len(_balancetes); t++ {

			SumValuesFathers(_balancetes[t].CodAccount, _balancetes)

			// fmt.Println(_balancetes[t].CodAccount, " - ", _balancetes[t].Description, _balancetes[t].DebitValue, _balancetes[t].CreditValue)
		}

	}
	return _balancetes
}
