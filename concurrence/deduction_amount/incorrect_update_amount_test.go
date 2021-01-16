// package goclub_book_concurrence_deduction_amount
//
// import (
// 	"context"
// 	"database/sql"
// 	_ "github.com/go-sql-driver/mysql"
// 	sq "github.com/goclub/sql"
// 	"github.com/pkg/errors"
// 	"log"
// 	"sync"
// 	"testing"
// )
//
// var db *sq.Database
// func init() {
// 	var err error
// 	// 演示场景下忽略 dbClose
// 	db, _, err = sq.Open("mysql", "root:somepass@(localhost:3306)/goclub_book") ; if err != nil {
// 		panic(err)
// 	}
// }
//
// type TableAmount struct {sq.WithoutSoftDelete}
// func (TableAmount) TableName() string { return "amount" }
//
// func QueryAmountByID(id uint64) (amount float64, err error) {
// 	// SELECT `amount` FROM `deduction_amount` WHERE `id` = ? LIMIT 1
// 	hasName, err := db.QueryRowScan(context.TODO(), sq.QB{
// 				Debug: true,
// 				Table: TableAmount{},
// 				Select: []sq.Column{"amount"},
// 				Where: sq.And("id", sq.Equal(id)),
// 	}, &amount)
// 	if err != nil { return }
// 	if !hasName {
// 		err = errors.New("not found data") ; return
// 	}
// 	return
// }
// func ResetAmountByID(id uint64, amount float64) {
//
// 	return
// }
//
//
// func TestIncorrectUpdateSetFixedAmount(t *testing.T) {
// 	wg := sync.WaitGroup{}
// 	name := "IncorrectUpdateSetFixedAmount"
// 	ResetAmountByName(name, 100)
// 	for i:=0;i<10;i++ {
// 		wg.Add(1)
// 		// 必须使用 routine 才能模拟并发
// 		go func() {
// 			amount := QueryAmountByName(name)
// 			updateAmount := amount - 1
// 			if updateAmount < 0 {
// 				log.Print("余额不够")
// 				return
// 			}
// 			resetAmountSQL := "UPDATE `deduction_amount` SET `amount` = ? WHERE `name` = ?"
// 			_, err := db.Exec(resetAmountSQL, updateAmount, name) ; if err != nil { panic(err) }
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }
//
// func TestIncorrectSubtractAmountOutOfRange(t *testing.T) {
// 	name := "IncorrectSubtractAmountOutOfRange"
// 	wg := sync.WaitGroup{}
// 	ResetAmountByName(name, 5)
// 	// 运行会出现 Out of range value for column 'amount' at row 1
// 	for i:=0;i<10;i++ {
// 		wg.Add(1)
// 		go func() {
// 			resetAmountSQL := "UPDATE `deduction_amount` SET `amount` = `amount` - ?   WHERE `name` = ?"
// 			_, err := db.Exec(resetAmountSQL, 1, name) ; if err != nil { panic(err) }
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }
//
// // compare and swap
// // 乐观锁：通过 Where 增加修改条件，并检查修改结果，但是 RowsAffected 可能不被支持
// func TestSubtractAmountSafeRange(t *testing.T) {
// 	name := "SubtractAmountSafeRange"
// 	wg := sync.WaitGroup{}
// 	ResetAmountByName(name, 5)
// 	subtractAmount := 1
// 	for i:=0;i<10;i++ {
// 		wg.Add(1)
// 		go func() {
// 			// 注意：
// 			// WHERE `amount` - ? >= 0 在遇到 unsigned 时会报错,使用 WHERE `amount` >= ? 可避免错误
// 			resetAmountSQL := "UPDATE `deduction_amount` SET `amount` = `amount` - ? WHERE `amount` >= ? AND `name` = ?"
// 			result, err := db.Exec(resetAmountSQL, subtractAmount, subtractAmount, name) ; if err != nil {
// 				panic(err)
// 			}
// 			// RowsAffected returns the number of rows affected by an
// 			// update, insert, or delete. Not every database or database
// 			// driver may support this.
// 			affected, err :=result.RowsAffected() ; if err != nil {
// 				panic(err)
// 			}
// 			if affected == 0 {
// 				log.Print("数据没有被修改")
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }
//
// func TestUpdateTransaction(t *testing.T) {
// 	name := "UpdateTransaction"
// 	wg := sync.WaitGroup{}
// 	ResetAmountByName(name, 5)
// 	var subtractAmount float64 = 1
// 	for i:=0;i<10;i++ {
// 		wg.Add(1)
// 		// i 如果不通过赋值传入函数会出现 print 的始终是10
// 		go func(i int) {
// 			log.Print("启动 routine: ", i)
// 			defer wg.Done()
// 			tx, err := db.BeginTxx(context.TODO(), &sql.TxOptions{sql.LevelRepeatableRead, false}) ; if err != nil {
// 				panic(err)
// 			}
// 			// 使用主键 where 条件能控制锁影响行最小
// 			queryAmountLockUpdateSQL := "SELECT `id`, `amount` FROM `deduction_amount` WHERE id = ? FOR UPDATE"
// 			row := tx.QueryRow(queryAmountLockUpdateSQL, 4)
// 			var id int
// 			var amount float64
// 			err = row.Scan(&id, &amount)
// 			has, err := CheckScanError(err) ; if err != nil {
// 				tx.Rollback()
// 				panic(err)
// 			}
// 			if !has {panic(errors.New("没有数据"))}
// 			updatedAmount := amount - subtractAmount
// 			if updatedAmount < 0 {
// 				log.Print("余额不够修改失败（routine: ", i, ")")
// 				tx.Rollback()
// 				return
// 			}
// 			// 保持惯例，在任何场景涉及到扣除数字的都使用 field = field - ? 的方式
// 			_, err = tx.Exec("UPDATE `deduction_amount` SET `amount` = `amount` - ? WHERE `id` = ?", subtractAmount, 4) ; if err != nil {
// 				tx.Rollback()
// 				panic(err)
// 			}
// 			log.Print("修改成功（routine: ", i, "）:", amount)
// 			tx.Commit()
// 		}(i)
// 	}
// 	wg.Wait()
// }
//
//
//
// /*
// 	不合适的乐观锁在遇到并发时执行结果大部分都会"失败"（严谨的方式应该使用 version 作为 compare 条件）
// */
// func TestIncorrectCompareAndSwap(t *testing.T) {
// 	name := "IncorrectCompareAndSwap"
// 	wg := sync.WaitGroup{}
// 	ResetAmountByName(name, 5)
// 	var subtractAmount float64 = 1
// 	for i:=0;i<10;i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			queryAmountSQL := "SELECT `id`, `amount` FROM `deduction_amount` WHERE `name` = ?"
// 			row := db.QueryRowx(queryAmountSQL, name)
// 			var id int
// 			var amount float64
// 			scanErr := row.Scan(&id, &amount)
// 			oldAmount := amount
// 			has, err := CheckScanError(scanErr) ; if err != nil {
// 				panic(err)
// 			}
// 			if !has { panic(errors.New("can not data:" + name)) }
// 			modifyAmount := amount - subtractAmount
// 			if modifyAmount < 0 {
// 				log.Print("余额不够！")
// 				return
// 			}
// 			updateSQL := "UPDATE `deduction_amount` SET `amount` = ? WHERE `id` = ? AND `amount` = ?"
// 			result, err := db.Exec(updateSQL, modifyAmount, id, oldAmount) ; if err != nil {
// 				panic(err)
// 			}
// 			rowsAffectedCount, err := result.RowsAffected() ; if err != nil {
// 				panic(err)
// 			}
// 			if rowsAffectedCount == 0 {
// 				log.Print("失败")
// 			} else {
// 				log.Print("成功")
// 			}
// 		}(i)
// 	}
// 	wg.Wait()
// }