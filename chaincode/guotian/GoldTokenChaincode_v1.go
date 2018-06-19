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
    ID           string  `json:"id"`               //ID （唯一ID）
    Name         string  `json:"name"`              //数字黄金名称: eg Simon Bucks
    Decimals     int     `json:"decimals"`            //小数点位数长度.
    Symbol       string  `json:"symbol"`            //标识: eg SBX
    Version      string  `json:"version"`           //版本信息 eg：1.0
    Owner        string  `json:"owner"`              //发起人ID地址
    TotalSupply  int     `json:"totalSupply"`         //数字黄金总量
    //balances map[Address] int
    //OwnerID        int  `json:"ownerId"`            //发行人ID
}

/**
* 用户数据结构
**/
type User struct {
	Address     string   `json:"address"`     //用户地址（唯一ID）
	Name        string   `json:"name"`        //用户名称
	Mobile      string   `json:"mobile"`      //手机号码
	Amount      int      `json:"amount"`      //账户余额（数字黄金持有数量）
	CardName    string   `json:"cardName"`    //身份证姓名
	CardNo      string   `json:"cardNo"`      //身份证号码
	IsCardOauth int      `json:"isCardOauth"` //身份身份证实名认证 0: 否 1：是
	Time        string   `json:"time"`        //创建时间
}


/**
* 中央银行
* 
* 负责数字黄金Token的发行、回收、销毁TOKEN等
* 只能和Bank进行交易流通TOKEN
*/
type CenterBank2 struct {
	Address     string  `json:"address"`     //中央银行地址（唯一ID）
	ID          string  `json:"id"`          //中央银行ID(唯一ID，不可修改)
	NameCN      string  `json:"nameCn"`      //中央银行名称(中文名称)
	NameEN      string  `json:"nameEn"`      //中央银行名称(英文名称)
	CompanyCode string  `json:"companyCode"` //企业社会信用统一代码
	Amount      int     `json:"amount"`      //账户余额（数字黄金持有数量）
	FrozenAmount  int     `json:"frozenAmount"`      //账户冻结余额（数字黄金持有数量）
	Time        string  `json:"time"`        //创建时间
}

/**
* 银行
*/
type Bank2 struct {
	Address     string  `json:"address"`     //银行地址（唯一ID）
	ID          string  `json:"id"`          //银行ID(唯一ID，不可修改)
	NameCN      string  `json:"nameCn"`      //银行名称(中文名称)
	NameEN      string  `json:"nameEn"`      //银行名称(英文名称)
	CompanyCode string  `json:"companyCode"` //企业社会信用统一代码
	Amount      int     `json:"amount"`      //账户余额（数字黄金持有数量）
	FrozenAmount  int     `json:"frozenAmount"`      //账户冻结余额（数字黄金持有数量）
	Time        string  `json:"time"`        //创建时间
}


/**
* 企业
*/
type Company2 struct {
	Address     string  `json:"address"`     //企业地址（唯一ID）
	ID          string  `json:"id"`          //企业ID(唯一ID，不可修改)
	NameCN      string  `json:"nameCn"`      //企业名称(中文名称)
	NameEN      string  `json:"nameEn"`      //企业名称(英文名称)
	CompanyCode string  `json:"companyCode"` //企业社会信用统一代码
	Amount      int     `json:"amount"`      //账户余额（数字黄金持有数量）
	FrozenAmount  int     `json:"frozenAmount"`      //账户冻结余额（数字黄金持有数量）
	Time        string  `json:"time"`        //创建时间

}

//中央银行
type CenterBank struct {
	Name        string `json:"name"`        //中央银行名称
	TotalNumber int    `json:"totalnumber"` //发行货币总数额
	RestNumber  int    `json:"restnumber"`  //账户余额
	ID          int    `json:"id"`          //中央银行ID
}

//银行
type Bank struct {
	Name        string `json:"name"`        //银行名称
	TotalNumber int    `json:"totalnumber"` //接收货币总数额
	RestNumber  int    `json:"fromtype"`    //账户余额
	ID          int    `json:"id"`          //银行ID
}

//企业
type Company struct {
	Name   string `json:"name"`   //企业名称
	Number int    `json:"number"` //账户余额
	ID     int    `json:"id"`     //企业ID

}

/**
* Token 新增、销毁交易
* 只有央行才可以新增、销毁TOKEN的数量
*/
type TokenIssueTran struct {
	Address   string `json:"address"`   //交易地址
	Type      int    `json:"type"`      //交易类型 1:新增 2:销毁
	Amount    int    `json:"amount"`    //数量
	Time      string `json:"time"`      //交易时间
    Remarke   string `json:"remarke"`   //备注说明
}


/**
* TOKEN 交易记录
*/
type Transaction2 struct {
	Address      string  `json:"address"`      //交易地址
	OrderNo      string  `json:"orderNo"`      //交易编号（订单交易号）
	FromRole     int     `json:"fromRole"`     //发送方角色 CenterBank:1, Bank:2, Company:3, User:4
	FromAddress  string  `json:"fromAddress"`  //发送方 ID
	ToRole       int     `json:"toRole"`       //接收方角色 CenterBank:1, Bank:2, Company:3, User:4
	ToAddress    string  `json:"toAddress"`    //接收方 ID
	Amount       int     `json:"amount"`       //交易数量
	Time         string  `json:"time"`         //交易时间
	Remarke      string  `json:"remarke"`      //备注说明
	
}

/**
* 交易内容
*/
type Transaction struct {
	FromType string `json:"fromtype"` //发送方角色 centerBank:0,Bank:1,Company:2
	FromID   int    `json:"fromid"`   //发送方 ID
	ToType   string `json:"totype"`   //接收方角色 Bank:1,Company:2
	ToID     int    `json:"toid"`     //接收方 ID
	Time     string `json:"time"`     //交易时间
	Number   int    `json:"number"`   //交易数额
	ID       int    `json:"id"`       //交易 ID
}

//0x2b38055e72da99f7ada2f09dd4e08951f5c8d52c984ecafcd9c4faee8a3ddf57
// 地址
var GoldToken_ADDRESS = string("0x9FE166aa9cF5BbFDBAf31e429E9923D994dB5199")
var CenterBank_ADDRESS = string("48d877acf2a04e63b5c2cdaffda97427")

/**
* 数字黄金对象实例
**/
var token GoldToken

/**
* 中央银行对象实例
**/
var center CenterBank

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
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	
	var token_name string        //数字黄金名称
    //var token_decimals uint8     //小数点位数长度
    var token_symbol string      //标识
    var token_version string      //版本
    var token_totalSupply int     //总量

	var err error

    var cbank CenterBank2        // 中央银行实例对象  
    
    if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// 初始化中央银行实例对象
	cbank.Address = CenterBank_ADDRESS
	cbank.ID = string("cb0")
	cbank.NameCN = string("国天黄金供应链（深圳）有限公司")
	cbank.NameEN = string("Guotian Gold Supply Chain（Shenzhen）Co.,Ltd.")
	cbank.Amount = 0
	cbank.FrozenAmount = 0
	cbank.CompanyCode = string("91320991056623231X")
	cbank.Time = string("2018-01-01 00:00:00")

    fmt.Printf("CenterBank Object cbank property Address = %d, ID = %d, NameCN=%d, NameEN=%d\n", cbank.Address, cbank.ID, cbank.NameCN, cbank.NameEN)

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
	token.ID = GoldToken_ADDRESS
	token.Name = token_name
	token.Decimals = 4
    token.Symbol = token_symbol
    token.Version = token_version
    token.TotalSupply = token_totalSupply
    token.Owner = CenterBank_ADDRESS

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
* 创建银行账号
*  
* 参数
* args[0] b_address 银行地址
* args[1] b_id 银行账号ID（唯一ID）
* args[2] b_name_cn 银行名称(中文)
* args[3] b_time 创建时间
*/
func (t *GoldTokenChaincode) CreateBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> CreateBank")

	var b_address string //银行地址
	var b_id string //银行账号ID（唯一ID）
    var b_name_cn string //银行名称(中文)
    //var b_company_code string //手机号码
    var b_time string //创建时间

    var err error
     
	var bank Bank2  //银行对象

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 参数设置
	b_address = args[0]
	b_id = args[1]
	b_name_cn = args[2]
	b_time = args[3]

    fmt.Printf(" b_address = %d, b_id  = %d, b_name_cn =%d, b_time =%d\n", b_address, b_id, b_name_cn, b_time)

    //cur_time := time.Now()
	
	// 初始化银行实例对象
	bank.Address = b_address
	bank.ID = b_id
	bank.NameCN = b_name_cn
	//bank.NameEN = string("Guotian Gold Supply Chain（Shenzhen）Co.,Ltd.")
	bank.Amount = 0
	bank.FrozenAmount = 0
	//bank.CompanyCode = b_company_code
	bank.Time = b_time
    //bank.Time = cur_time.String() //获取当前时间


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
* 创建企业账号
*  
* 参数
* args[0] b_address 企业地址
* args[1] b_id 企业账号ID（唯一ID）
* args[2] b_name_cn 企业名称(中文)
* args[3] b_time 创建时间
*/
func (t *GoldTokenChaincode) CreateCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> CreateCompany")

	var b_address string //企业地址
	var b_id string //企业账号ID（唯一ID）
    var b_name_cn string //企业名称(中文)
    //var b_company_code string //手机号码
    var b_time string //创建时间

    var err error
     
	var company Company2  //用户对象

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 参数设置
	b_address = args[0]
	b_id = args[1]
	b_name_cn = args[2]
	b_time = args[3]

    fmt.Printf(" b_address = %d, b_id  = %d, b_name_cn =%d, b_time =%d\n", b_address, b_id, b_name_cn, b_time)

    //cur_time := time.Now()
	
	// 初始化中央银行实例对象
	company.Address = b_address
	company.ID = b_id
	company.NameCN = b_name_cn
	//company.NameEN = string("Guotian Gold Supply Chain（Shenzhen）Co.,Ltd.")
	company.Amount = 0
	company.FrozenAmount = 0
	//company.CompanyCode = b_company_code
	company.Time = b_time
    //company.Time = cur_time.String() //获取当前时间


	jsons, errs := json.Marshal(company) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(company.Address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    return shim.Success(jsons)
}

/**
* 创建用户
* args[0] Address 用户地址
* args[1] Name 用户名字
* args[2] Mobile 手机号码
* args[3] Time 注册时间
**/
func (t *GoldTokenChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode createUser")

    var u_address string //用户地址
    var u_name string //用户姓名
    var u_mobile string //手机号码
    var u_time string //创建时间

    var err error
     
	var user User  //用户对象

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 参数设置
	u_address = args[0]
	u_name = args[1]
	u_mobile = args[2]
	u_time = args[3]

    fmt.Printf(" u_address = %d, u_name  = %d,u_mobile =%d, u_time =%d\n", u_address, u_name, u_mobile, u_time)

    cur_time := time.Now()
	
	//初始化user对象
    user.Address = u_address
    user.Name = u_name
    user.Mobile = u_mobile
    user.Amount = 0
    user.IsCardOauth = 0 //没有身份证实名认证
    user.Time = cur_time.String() //获取当前时间


	jsons, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(u_address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    //jsonResp := "{\"Address00000\":\"1111111\",\"Name\":\"888888\"}"
    
    return shim.Success(jsons)
}

/**
* 用户身份证认证
* args[0] Address 用户地址
**/
func (t *GoldTokenChaincode) oauthUserCard(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode oauthUserCard")

    var u_address   string //用户地址
    var u_cardName  string //身份证姓名
    var u_cardNo    string //身份证号码
    
    var err error
     
	var user User  //用户对象

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// 参数解析
	u_address = args[0]
	u_cardName = args[1]
	u_cardNo = args[2]
	
    fmt.Printf(" u_address = %d, u_cardName  = %d, u_cardNo =%d\n", u_address, u_cardName, u_cardNo)

    UserBytes, erro := stub.GetState(u_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(UserBytes, &user)
	user.CardName = u_cardName
	user.CardNo = u_cardNo
	user.IsCardOauth = 1 //已经身份证实名认证
    
	jsons, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(u_address, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

    return shim.Success(jsons)
	//return shim.Success(nil)
}

/**
* 
* 
* 新增TOKEN数量
* 参数:
* args[0] address 交易地址
* args[1] amount  新增数量
* args[2] tx_time 交易地址
* args[3] remarke 交易地址
**/
func (t *GoldTokenChaincode) IssueCoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> IssueCoin")

    var address string     //交易地址
	var amount int        //新增金额
	var tx_time string     //交易时间
	var remarke string      //备注说明
    
    var trans TokenIssueTran //交易过程

    var cbank CenterBank2    // 中央银行实例对象
    var goldToken GoldToken  // 数字黄金token对象

	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	address = args[0]
	amount, err = strconv.Atoi(args[1])
    if err != nil {
		return shim.Error("Expecting integer value for asset holding：amount ")
	}
    tx_time = args[2]
    remarke = args[3]


	fmt.Printf("  address  = %d , amount = %d \n", address, amount)



    //初始化对象
    trans.Address = address //交易地址
    trans.Type = 1          //新增:1 销毁:2 
    trans.Amount = amount   //金额
    //cur_time := time.Now()
	//trans.Time = cur_time.String() //创建时间
	trans.Time = tx_time //创建时间
	trans.Remarke = remarke  //备注说明

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
	var tx_time string     //交易时间
	var remarke string      //备注说明
    
    var trans TokenIssueTran //交易过程

    var cbank CenterBank2    // 中央银行实例对象
    var goldToken GoldToken  // 数字黄金token对象

	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	address = args[0]
	amount, err = strconv.Atoi(args[1])
    if err != nil {
		return shim.Error("Expecting integer value for asset holding：amount ")
	}
    tx_time = args[2]
    remarke = args[3]


	fmt.Printf("  address  = %d , amount = %d \n", address, amount)



    //初始化对象
    trans.Address = address //交易地址
    trans.Type = 2          //新增:1 销毁:2 
    trans.Amount = amount   //金额
    //cur_time := time.Now()
	//trans.Time = cur_time.String() //创建时间
	trans.Time = tx_time //创建时间
	trans.Remarke = remarke  //备注说明

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
* 过程: [中央银行] --> [商业银行]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransCb2Bank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransCb2Bank")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var cbank CenterBank2
    var bank  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 1
    trans.FromAddress = from_address
    trans.ToRole = 2
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

	CenterBankBytes, erro := stub.GetState(from_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CenterBankBytes, &cbank)
	cbank.Amount = cbank.Amount - amount

	jsons_cbank, errs := json.Marshal(cbank) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(cbank.Address, jsons_cbank)

	BankBytes, erro1 := stub.GetState(to_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount + amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

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

/**
* 转账功能
* 
* 过程: [商业银行] --> [中央银行]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransBank2Cb(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2Cb")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var cbank CenterBank2
    var bank  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 2
    trans.FromAddress = from_address
    trans.ToRole = 1
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

    BankBytes, erro1 := stub.GetState(from_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount - amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

	CenterBankBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CenterBankBytes, &cbank)
	cbank.Amount = cbank.Amount + amount

	jsons_cbank, errs := json.Marshal(cbank) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(cbank.Address, jsons_cbank)


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

/**
* 转账功能
* 
* 过程: [商业银行] --> [商业银行]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransBank2Bank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2Bank")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var bankFrom Bank2
    var bankTo  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 2
    trans.FromAddress = from_address
    trans.ToRole = 2
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

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


/**
* 转账功能
* 
* 过程: [商业银行] --> [企业商家]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransBank2Cp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2Cp")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var company Company2
    var bank  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 2
    trans.FromAddress = from_address
    trans.ToRole = 3
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

    BankBytes, erro1 := stub.GetState(from_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount - amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

	CompanyBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CompanyBytes, &company)
	company.Amount = company.Amount + amount

	jsons_company, errs := json.Marshal(company) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(company.Address, jsons_company)


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


/**
* 转账功能
* 
* 过程: [企业商家] --> [商业银行]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransCp2Bank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2Cp")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var company Company2
    var bank  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 3
    trans.FromAddress = from_address
    trans.ToRole = 2
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke


    CompanyBytes, erro := stub.GetState(from_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CompanyBytes, &company)
	company.Amount = company.Amount - amount

	jsons_company, errs := json.Marshal(company) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(company.Address, jsons_company)

    BankBytes, erro1 := stub.GetState(to_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount + amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

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


/**
* 转账功能
* 
* 过程: [商业银行] --> [普通用户]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransBank2User(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2User")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var user User
    var bank  Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 2
    trans.FromAddress = from_address
    trans.ToRole = 4
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

    BankBytes, erro1 := stub.GetState(from_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount - amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

	UserBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(UserBytes, &user)
	user.Amount = user.Amount + amount

	jsons_user, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(user.Address, jsons_user)


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


/**
* 转账功能
* 
* 过程: [普通用户] --> [商业银行]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransUser2Bank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransBank2User")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var user User
    var bank Bank2

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 4
    trans.FromAddress = from_address
    trans.ToRole = 2
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

    UserBytes, erro := stub.GetState(from_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(UserBytes, &user)
	user.Amount = user.Amount - amount

	jsons_user, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(user.Address, jsons_user)

    BankBytes, erro1 := stub.GetState(to_address)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankBytes, &bank)
	bank.Amount = bank.Amount + amount
	jsons_bank, errs2 := json.Marshal(bank) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(bank.Address, jsons_bank)

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


/**
* 转账功能
* 
* 过程: [企业商家] --> [普通用户]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransCp2User(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransCp2User")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var company Company2
    var user User

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 3
    trans.FromAddress = from_address
    trans.ToRole = 4
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke


    CompanyBytes, erro := stub.GetState(from_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CompanyBytes, &company)
	company.Amount = company.Amount - amount

	jsons_company, errs := json.Marshal(company) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(company.Address, jsons_company)

    UserBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(UserBytes, &user)
	user.Amount = user.Amount + amount

	jsons_user, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(user.Address, jsons_user)

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


/**
* 转账功能
* 
* 过程: [普通用户] --> [企业商家]
* 参数:
* args[0] address 交易地址
* args[1] orderNo 订单编号
* args[2] from_address  转账方地址
* args[3] to_address 接受方地址
* args[4] amount 转账数量
* args[5] tx_time 交易时间
* args[6] remarke 备注说明
**/
func (t *GoldTokenChaincode) TransUser2Cp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode --> TransUser2Cp")
    // 参数定义
    var address string  //交易地址
    var orderNo string  // 订单编号
    var from_address string  //转账方地址
    var to_address string    //接受方地址
    var amount int   //转账数量
    var tx_time string //交易时间
    var remarke string //备注说明
	
	var err error
    
    var trans Transaction2 //交易对象
    var company Company2
    var user User

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
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
	tx_time = args[5]
	remarke = args[6]

    trans.Address = address
    trans.OrderNo = orderNo
    trans.FromRole = 4
    trans.FromAddress = from_address
    trans.ToRole = 3
    trans.ToAddress = to_address
    trans.Amount = amount
    trans.Time = tx_time
    trans.Remarke = remarke

    UserBytes, erro := stub.GetState(from_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(UserBytes, &user)
	user.Amount = user.Amount - amount

	jsons_user, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(user.Address, jsons_user)


    CompanyBytes, erro := stub.GetState(to_address)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(CompanyBytes, &company)
	company.Amount = company.Amount + amount

	jsons_company, errs := json.Marshal(company) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	// Write the state to the ledger
	err = stub.PutState(company.Address, jsons_company)


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

//商业银行转账到企业  issueCoinToCp
func (t *GoldTokenChaincode) issueCoinToCp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ex02 IssueCoin")

	var Number int             // 数量
	var From_ID int            // 商业银行ID
	var To_ID int              //接收方ID
	var ID int                 //交易ID
	var bank_to_cp Transaction //交易过程
	var bankFrom Bank          //商业银行
	var cpTo Company           //企业
	var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode

	From_ID, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：From_ID ")
	}
	Number, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：Number ")
	}
	To_ID, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：To_ID  ")
	}

	ID, err = strconv.Atoi(args[3])

	if err != nil {
		return shim.Error("Expecting integer value for asset holding：ID_trans ")
	}

	fmt.Printf("  Number  = %d ,To_ID =%d , ID_trans=%d\n", Number, To_ID, ID)

	bank_to_cp.FromType = "1"
	bank_to_cp.FromID = From_ID
	bank_to_cp.ToType = "2"
	bank_to_cp.ToID = To_ID

	cur_time := time.Now()
	bank_to_cp.Time = cur_time.String()

	bank_to_cp.Number = Number
	bank_to_cp.ID = ID

	BankFromBytes, erro := stub.GetState(args[0])
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(BankFromBytes, &bankFrom)
	bankFrom.RestNumber = bankFrom.RestNumber - Number

	jsons_bank, errs := json.Marshal(bankFrom) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	bankFromID_string := strconv.Itoa(bankFrom.ID)

	// Write the state to the ledger
	err = stub.PutState(bankFromID_string, jsons_bank)

	companyToBytes, erro1 := stub.GetState(args[1])
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(companyToBytes, &cpTo)
	cpTo.Number = cpTo.Number + Number

	jsons_cp, errs2 := json.Marshal(cpTo) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	cpToID_string := strconv.Itoa(cpTo.ID)
	// Write the state to the ledger
	err = stub.PutState(cpToID_string, jsons_cp)

	jsons, errs3 := json.Marshal(bank_to_cp) //转换成JSON返回的是byte[]
	if errs3 != nil {
		return shim.Error(errs3.Error())
	}
	ID_string := strconv.Itoa(ID)
	// Write the state to the ledger
	err = stub.PutState(ID_string, jsons)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//getBanks
func (t *GoldTokenChaincode) getBanks(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ex02 getBanks")

	var Bank_ID string // 商业银行ID
	var bank_info Bank
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode

	Bank_ID = args[0]

	BankInfo, erro := stub.GetState(Bank_ID)
	if erro != nil {
		return shim.Error(erro.Error())
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(BankInfo, &bank_info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("  BankInfo  = %d  \n", BankInfo)

	return shim.Success(nil)
}

//getCompanys
func (t *GoldTokenChaincode) getCompanys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ex02 getCompanys")

	var CP_ID string // 企业ID
	var company_info Company
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode

	CP_ID = args[0]

	company_info_bytes, erro := stub.GetState(CP_ID)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct

	err = json.Unmarshal(company_info_bytes, &company_info)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("  BankInfo  = %d  \n", company_info_bytes)

	return shim.Success(nil)
}

//getTransactions
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

//getCenterBank
func (t *GoldTokenChaincode) getCenterBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("ex02 getCenterBank")

	var Center_ID string // 企业ID
	var center_info CenterBank
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode

	Center_ID = args[0]

	center_info_bytes, erro := stub.GetState(Center_ID)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct

	err = json.Unmarshal(center_info_bytes, &center_info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("  center_info_bytes  = %d  \n", center_info_bytes)

	return shim.Success(nil)
}

//transfer
func (t *GoldTokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("GoldTokenChaincode -> getCenterBank")

	var From_ID int // 转账方ID
	var To_ID int   //接收方ID
	var number int  //转账金额
	var fromCP Company
	var toCP Company
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// Initialize the chaincode

	From_ID, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：From_ID  ")
	}
	To_ID, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：To_ID  ")
	}
	number, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding：number ")
	}

	fromID_string := strconv.Itoa(From_ID)
	from_cp_info_bytes, erro := stub.GetState(fromID_string)
	if erro != nil {
		return shim.Error(erro.Error())
	}

	//将byte的结果转换成struct

	err = json.Unmarshal(from_cp_info_bytes, &fromCP)

	fmt.Printf("  from_cp_info_bytes  = %d  \n", from_cp_info_bytes)

	To_ID_string := strconv.Itoa(To_ID)
	to_cp_info_bytes, erro1 := stub.GetState(To_ID_string)
	if erro1 != nil {
		return shim.Error(erro1.Error())
	}

	//将byte的结果转换成struct

	err = json.Unmarshal(to_cp_info_bytes, &toCP)

	fmt.Printf("  to_cp_info_bytes  = %d  \n", to_cp_info_bytes)

	from_cp_old_num := fromCP.Number
	if from_cp_old_num <= number {
		return shim.Error("money no enough")
	}

	fromCP.Number = from_cp_old_num - number

	to_cp_old_num := toCP.Number
	toCP.Number = to_cp_old_num + number

	jsons_from, errs := json.Marshal(fromCP) //转换成JSON返回的是byte[]
	if errs != nil {
		return shim.Error(errs.Error())
	}
	fromCPID_string := strconv.Itoa(fromCP.ID)
	// Write the state to the ledger
	err = stub.PutState(fromCPID_string, jsons_from)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsons_to, errs2 := json.Marshal(toCP) //转换成JSON返回的是byte[]
	if errs2 != nil {
		return shim.Error(errs2.Error())
	}
	toCPID_string := strconv.Itoa(toCP.ID)
	// Write the state to the ledger
	err = stub.PutState(toCPID_string, jsons_to)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf(" transfer success \n")
	return shim.Success(nil)
}

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
	fmt.Println("ex02 Invoke")
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
	} else if function == "CreateBank" {
		// the old "Query" is now implemtned in invoke
		return t.CreateBank(stub, args)
	} else if function == "CreateCompany" {
		// the old "Query" is now implemtned in invoke
		return t.CreateCompany(stub, args)
	} else if function == "getBanks" {
		// the old "Query" is now implemtned in invoke
		return t.getBanks(stub, args)
	} else if function == "getCenterBank" {
		// the old "Query" is now implemtned in invoke
		return t.getCenterBank(stub, args)
	} else if function == "getCompanys" {
		// the old "Query" is now implemtned in invoke
		return t.getCompanys(stub, args)
	} else if function == "getTransactions" {
		// the old "Query" is now implemtned in invoke
		return t.getTransactions(stub, args)
	} else if function == "IssueCoin" {
		// 新增Token数量
		return t.IssueCoin(stub, args)
	} else if function == "DestroyCoin" {
		// 销毁Token数量
		return t.DestroyCoin(stub, args)
	} else if function == "issueCoinToBank" {
		// the old "Query" is now implemtned in invoke
		return t.issueCoinToBank(stub, args)
	} else if function == "issueCoinToCp" {
		// the old "Query" is now implemtned in invoke
		return t.issueCoinToCp(stub, args)
	} else if function == "transfer" {
		// the old "Query" is now implemtned in invoke
		return t.transfer(stub, args)
	} else if function == "createUser" {
		// 创建用户
		return t.createUser(stub, args)
	} else if function == "oauthUserCard" {
		// 用户身份认证
		return t.oauthUserCard(stub, args)
	} else if function == "TransCb2Bank" {
		// 转账交易 中央银行 -》商业银行
		return t.TransCb2Bank(stub, args)
	} else if function == "TransBank2Cb" {
		// 转账交易 商业银行 -》中央银行
		return t.TransBank2Cb(stub, args)
	} else if function == "TransBank2Bank" {
		// 转账交易 商业银行 -》商业银行
		return t.TransBank2Bank(stub, args)
	} else if function == "TransBank2Cp" {
		// 转账交易 商业银行 -》企业商家
		return t.TransBank2Cp(stub, args)
	} else if function == "TransCp2Bank" {
		// 转账交易 企业商家 -》商业银行
		return t.TransCp2Bank(stub, args)
	} else if function == "TransBank2User" {
		// 转账交易 商业银行 -》普通用户
		return t.TransBank2User(stub, args)
	} else if function == "TransUser2Bank" {
		// 转账交易 普通用户 -》商业银行
		return t.TransUser2Bank(stub, args)
	} else if function == "TransCp2User" {
		// 转账交易 企业商家 -》普通用户
		return t.TransCp2User(stub, args)
	} else if function == "TransUser2Cp" {
		// 转账交易 普通用户 -》企业商家
		return t.TransUser2Cp(stub, args)
	}

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
