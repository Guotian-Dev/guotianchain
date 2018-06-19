package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/**
*   数字黄金链码
*/
type GoldTokenChaincode struct {

}

/**
* 地址结构
*/
//type Address string

/**
* 数字黄金
* 注意：结构体转换成json 字段必须大写开头
*/
type GoldToken struct {
    Address           string  `json:"address"`               //ID （唯一ID）
    Name         string  `json:"name"`              //数字黄金名称: eg Simon Bucks
    Decimals     int     `json:"decimals"`            //小数点位数长度.
    Symbol       string  `json:"symbol"`            //标识: eg SBX
    Version      string  `json:"version"`           //版本信息 eg：1.0
    Owner        string  `json:"owner"`              //发起人ID地址
    TotalSupply  int     `json:"totalSupply"`         //数字黄金总量
    //balances map[Address] int
    //OwnerID        int  `json:"ownerId"`            //发行人ID
    Time        string  `json:"time"`          //创建时间
}


/**
* 账户
*/
type Account struct {
	Address     string  `json:"address"`       //账户地址（唯一ID）
	UserNo      string  `json:"userNo"`        //账户ID(唯一ID，不可修改)
	UserName    string  `json:"userName"`      //账户名称
	CardNo      string  `json:"cardNo"`        //身份证号码
	CardName    string  `json:"cardName"`      //身份证名称
	IsCardAuth  int     `json:"isCardAuth"`    //是否实名制认证（1-已认证 2-未认证）
	CompanyCode string  `json:"companyCode"`   //统一社会信用代码
	CompanyName string  `json:"companyName"`   //企业名称
	IsCompanyAuth  int  `json:"isCompanyAuth"` //是否企业认证（1-已认证 2-未认证）
	Amount      int     `json:"amount"`        //账户余额（持有Token数量）
	RoleCode    string  `json:"roleCode"`      //账户类型代码（发行商-10 运营商-11 企业-12 普通-13）
	Time        string  `json:"time"`          //创建时间
}


/**
* TOKEN 交易记录
*/
type Transaction struct {
	Address      string  `json:"address"`      //交易地址
	OrderNo      string  `json:"orderNo"`      //交易编号（订单交易号）
	FromAddress  string  `json:"fromAddress"`  //发送方 ID
	ToAddress    string  `json:"toAddress"`    //接收方 ID
	TypeCode     string  `json:"typeCode"`     //交易类型（增发Token-20 销毁Token-21 转账-22） 
	Amount       int     `json:"amount"`       //交易数量
	Time         string  `json:"time"`         //交易时间
	Remarke      string  `json:"remarke"`      //备注说明
	
}


//0x2b38055e72da99f7ada2f09dd4e08951f5c8d52c984ecafcd9c4faee8a3ddf57
// 地址
var GoldToken_ADDRESS = string("0x9FE166aa9cF5BbFDBAf31e429E9923D994dB5199")
var CenterBank_ADDRESS = string("48d877acf2a04e63b5c2cdaffda97427")

/**
* Token对象实例
**/
var token GoldToken


/**
* 链码初始化函数
*
* 参数说明
* args[0] 中央银行名称
* args[1] 发行货币总数额
* args[2] 账户余额
* args[3] 中央银行ID
*
* args[4] 数字黄金名称 
* args[5] 标示
* args[6] 版本
* args[7] 总量
* args[8] 总量
* 
**/
func (t *GoldTokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("GoldTokenChaincode --> Init")
	_, args := stub.GetFunctionAndParameters()
	
	var token_name string        //数字黄金名称
    //var token_decimals uint8     //小数点位数长度
    var token_symbol string      //标识
    var token_version string      //版本
    var token_totalSupply int     //总量

	var err error

    var cbank Account        // 发行商实例对象  
    
    if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// 初始化中央银行实例对象
	cbank.Address = CenterBank_ADDRESS
	cbank.UserNo = string("20180518000")
	cbank.UserName = string("国天黄金供应链（深圳）有限公司")
	cbank.RoleCode = string("10")
	cbank.Amount = 0
	cbank.IsCardAuth = 2
	cbank.IsCompanyAuth = 2
	cbank.CompanyCode = string("91320991056623231X")
	cbank.Time = string("2018-01-01 00:00:00")


    //fmt.Printf("CenterBank Object cbank property Address = %d, ID = %d, NameCN=%d, NameEN=%d\n", cbank.Address, cbank.ID, cbank.NameCN, cbank.NameEN)

    jsons_cb, errs_cb := json.Marshal(cbank) //转换成JSON返回的是byte[]

	if errs_cb != nil {
		return shim.Error(errs_cb.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(cbank.Address, jsons_cb)
	if err != nil {
		return shim.Error(err.Error())
	}


	// 初始化数字黄金对象
    token_name = args[4]
    //token_decimals = 4
    token_symbol = args[5]
    token_version = args[6]
    //token_id = GoldToken_ADDRESS

    //token_totalSupply, err = strconv.Atoi(args[7])
	// if err != nil {
	// 	return shim.Error("Expecting integer value for asset holding：totalSupply")
	// }

	token_totalSupply = 0
	
	//token对象设置
	token.Address = GoldToken_ADDRESS
	token.Name = token_name
	token.Decimals = 4
    token.Symbol = token_symbol
    token.Version = token_version
    token.TotalSupply = token_totalSupply
    token.Owner = CenterBank_ADDRESS

    cur_time := time.Now()
    token.Time = cur_time.String() //获取当前时间
    
    //token
    jsons_token, errs2 := json.Marshal(token) //转换成JSON返回的是byte[]

    fmt.Printf("id = %d, token = %d, jsons_token=%d\n", GoldToken_ADDRESS, token, jsons_token)


	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(GoldToken_ADDRESS, jsons_token)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf(" init success \n")
	return shim.Success(nil)
}

/**
* 创建账户
*  
* 参数
* args[0] user_address 账户地址
* args[1] user_no 账户号
* args[2] user_name 账户名称
* args[3] user_name 账户名称
*/
func (t *GoldTokenChaincode) CreateAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> CreateAccount")

	var user_address string //账户地址
	var user_no string      //账户号
    var user_name string    //账户名称
    var role_code string    //账户类型代码（央行-10 银行-11 企业-12 普通-13）
    
    var err error
     
	var bank Account  //账户对象

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 参数设置
	user_address = args[0]
	user_no = args[1]
	user_name = args[2]
	role_code = args[3]

    fmt.Printf("user_address = %d, user_no  = %d, user_name =%d, role_code =%d\n", user_address, user_no, user_name, role_code)

    cur_time := time.Now()
	
	// 初始化银行实例对象
	bank.Address = user_address
	bank.UserNo = user_no
	bank.UserName = user_name
	//bank.NameEN = string("Guotian Gold Supply Chain（Shenzhen）Co.,Ltd.")
	bank.Amount = 0
	bank.IsCardAuth = 2
	bank.IsCompanyAuth = 2
	bank.RoleCode = role_code
	//bank.Amount = 0
	//bank.CompanyCode = b_company_code
	//bank.Time = b_time
    bank.Time = cur_time.String() //获取当前时间


	jsons, errs := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    return shim.Success(jsons)
}


/**
* 
* 
* 新增TOKEN数量
* 参数:
* args[0] address 交易地址
* args[1] amount  新增数量
* args[2] tx_time 交易时间
* args[3] remarke 备注说明
**/
func (t *GoldTokenChaincode) IssueCoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> IssueCoin")

    var address string     //交易地址
	var amount int        //新增金额
	var orderNo string     //交易时间
	var remarke string      //备注说明
    
    var trans Transaction //交易过程

    var cbank Account    // 中央银行实例对象
    var goldToken GoldToken  // 数字黄金token对象

	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	address = args[0]
	orderNo = args[1]
	amount, err = strconv.Atoi(args[2])
    if err != nil {
		return shim.Error("Expecting integer value for asset holding：amount ")
	}
    remarke = args[3]


	fmt.Printf("  address  = %d , amount = %d \n", address, amount)


	cur_time := time.Now()

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromAddress = CenterBank_ADDRESS
    trans.ToAddress = CenterBank_ADDRESS
    trans.TypeCode = string("20")
    trans.Amount = amount
    trans.Remarke = remarke
    trans.Time = cur_time.String() //获取当前时间

    CenterBankBytes, erro := stub.GetState(CenterBank_ADDRESS)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CenterBankBytes, &cbank)
	cbank.Amount = cbank.Amount + amount 


	GoldTokenBytes, erro := stub.GetState(GoldToken_ADDRESS)
	if erro != nil {
		return shim.Error(erro.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(GoldTokenBytes, &goldToken)
	goldToken.TotalSupply = goldToken.TotalSupply + amount 
	
	
	//保存交易
	jsons, errs := json.Marshal(trans) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(trans.Address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    //保存中央银行
	jsons_cbank, errs2 := json.Marshal(cbank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(CenterBank_ADDRESS, jsons_cbank)
	if err != nil {
		return shim.Error(err.Error())
	}

	//保存GOLDTOKEN
	jsons_goldtoken, errs3 := json.Marshal(goldToken) //转换成JSON返回的是byte[]
	if errs3 != nil {
		return shim.Error(errs3.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(GoldToken_ADDRESS, jsons_goldtoken)
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Printf(" IssueCoin success \n")
    
    return shim.Success(jsons)
	//return shim.Success(nil)
}


/**
* 
* 
* 销毁TOKEN数量
* 参数:
* args[0] address 交易地址
* args[1] amount  销毁数量
* args[2] tx_time 交易地址
* args[3] remarke 交易地址
**/
func (t *GoldTokenChaincode) DestroyCoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> DestroyCoin")

    var address string     //交易地址
	var amount int        //新增金额
	var orderNo string     //交易单号
	var remarke string      //备注说明
    
    var trans Transaction //交易过程

    var cbank Account    // 中央银行实例对象
    var goldToken GoldToken  // 数字黄金token对象

	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	address = args[0]
	orderNo = args[1]
	amount, err = strconv.Atoi(args[2])
    if err != nil {
		return shim.Error("Expecting integer value for asset holding：amount ")
	}
    
    remarke = args[3]


	fmt.Printf("  address  = %d , amount = %d \n", address, amount)



    cur_time := time.Now()

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromAddress = CenterBank_ADDRESS
    trans.ToAddress = CenterBank_ADDRESS
    trans.TypeCode = string("21")
    trans.Amount = amount
    trans.Remarke = remarke
    trans.Time = cur_time.String() //获取当前时间

    CenterBankBytes, erro := stub.GetState(CenterBank_ADDRESS)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CenterBankBytes, &cbank)
	cbank.Amount = cbank.Amount - amount 


	GoldTokenBytes, erro := stub.GetState(GoldToken_ADDRESS)
	if erro != nil {
		return shim.Error(erro.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(GoldTokenBytes, &goldToken)
	goldToken.TotalSupply = goldToken.TotalSupply - amount 
	
	
	//保存交易
	jsons, errs := json.Marshal(trans) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(trans.Address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    //保存中央银行
	jsons_cbank, errs2 := json.Marshal(cbank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(CenterBank_ADDRESS, jsons_cbank)
	if err != nil {
		return shim.Error(err.Error())
	}

	//保存GOLDTOKEN
	jsons_goldtoken, errs3 := json.Marshal(goldToken) //转换成JSON返回的是byte[]
	if errs3 != nil {
		return shim.Error(errs3.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(GoldToken_ADDRESS, jsons_goldtoken)
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Printf(" IssueCoin success \n")
    
    return shim.Success(jsons)
	//return shim.Success(nil)
}


//----------------------------------------------------------------------------//

/**
* 转账功能
* 
* 过程: [账户] --> [账户]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] remarke 备注说明
**/
func (t *GoldTokenChaincode) Transaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> Transaction")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    //var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction //交易对象
    var bankFrom Account
    var bankTo  Account

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

    //参数设置
	address = args[0]
	orderNo = args[1]
	from_address = args[2]
	to_address = args[3]
	amount, err = strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：amount ")
	}
	//tx_time = args[5]
	remarke = args[5]

    cur_time := time.Now()

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromAddress = from_address
    trans.ToAddress = to_address
    trans.TypeCode = string("22")
    trans.Amount = amount
    trans.Remarke = remarke
    trans.Time = cur_time.String() //获取当前时间

    BankFromBytes, erro1 := stub.GetState(from_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankFromBytes, &bankFrom)
	bankFrom.Amount = bankFrom.Amount - amount
	jsons_bank_from, errs2 := json.Marshal(bankFrom) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bankFrom.Address, jsons_bank_from)

	BankToBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(BankToBytes, &bankTo)
	bankTo.Amount = bankTo.Amount + amount

	jsons_bank_to, errs := json.Marshal(bankTo) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bankTo.Address, jsons_bank_to)


	jsons, errs3 := json.Marshal(trans) //转换成JSON返回的是byte[]
	if errs3 != nil {
		return shim.Error(errs3.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(trans.Address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsons)
	//return shim.Success(nil)
}


//----------------------------------------------------------------------------//


//issueCoinToBank  发行货币至商业银行
/*
func (t *GoldTokenChaincode) issueCoinToBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ex02 IssueCoin")

	var Number int                // 发行的数量
	var To_ID int                 //接收方ID
	var ID_trans int              //交易ID
	var trans_to_bank Transaction //交易过程
	var toBank Bank               //商业银行
	var err error
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode

	Number, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：Number ")
	}
	To_ID, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：To_ID  ")
	}

	ID_trans, err = strconv.Atoi(args[2])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding：ID_trans ")
	}

	fmt.Printf("  Number  = %d ,To_ID =%d , ID_trans=%d\n", Number, To_ID, ID_trans)

	trans_to_bank.FromType = "0"
	trans_to_bank.FromID = 0
	trans_to_bank.ToType = "1"
	trans_to_bank.ToID = To_ID

	cur_time := time.Now()

	trans_to_bank.Time = cur_time.String()

	trans_to_bank.Number = Number
	trans_to_bank.ID = ID_trans

	center.RestNumber = center.RestNumber - Number

	toBankInfo, erro := stub.GetState(args[0])
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(toBankInfo, &toBank)
	toBank.TotalNumber = Number
	toBank.RestNumber = toBank.RestNumber + Number

	fmt.Printf("  toBankInfo  = %d  \n", toBankInfo)

	jsons, errs := json.Marshal(trans_to_bank) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	ID_trans_string := strconv.Itoa(ID_trans)
	// Write the state to the ledger
	err = stub.PutState(ID_trans_string, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsons_toBank, errs2 := json.Marshal(toBank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	toBankID_string := strconv.Itoa(toBank.ID)
	// Write the state to the ledger
	err = stub.PutState(toBankID_string, jsons_toBank)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsons_center, errs3 := json.Marshal(center) //转换成JSON返回的是byte[]
	if errs3 != nil {
		return shim.Error(errs3.Error())
	}
	centerID_string := strconv.Itoa(center.ID)
	// Write the state to the ledger
	err = stub.PutState(centerID_string, jsons_center)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("  issueCoinToBank success \n")
	return shim.Success(nil)
}
*/


//getTransactions
/*
func (t *GoldTokenChaincode) getTransactions(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> getTransactions")

	var trans_ID string // 企业ID
	var trans_info Transaction
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode

	trans_ID = args[0]

	trans_info_bytes, erro := stub.GetState(trans_ID)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct

	err = json.Unmarshal(trans_info_bytes, &trans_info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("  trans_info_bytes  = %d  \n", trans_info_bytes)

	return shim.Success(nil)
}
*/



// Deletes an entity from state
func (t *GoldTokenChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}
func (t *GoldTokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("GoldTokenChaincode --> Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "CreateAccount" {
		// 创建账户
		return t.CreateAccount(stub, args)
	} else if function == "Transaction" {
		// 转账交易
		return t.Transaction(stub, args)
	} else if function == "IssueCoin" {
		// 新增Token数量
		return t.IssueCoin(stub, args)
	} else if function == "DestroyCoin" {
		// 销毁Token数量
		return t.DestroyCoin(stub, args)
	}
	/* else if function == "issueCoinToBank" {
		// the old "Query" is now implemtned in invoke
		return t.issueCoinToBank(stub, args)
	} */

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *GoldTokenChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

// Deletes an entity from state

// query callback representing the query of a chaincode
func (t *GoldTokenChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, erro := stub.GetState(A)
	if erro != nil {
		return shim.Error(erro.Error())
	}
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(GoldTokenChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
