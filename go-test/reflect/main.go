package main
import (
"reflect"
"fmt"
)
type ControllerInterface interface {
	Init(action string, method string)
}
type Controller struct {
	Action string
	Method string
	Tag string `json:"tag"`
}
func (c *Controller) Init(action string, method string){
	c.Action = action
	c.Method = method
	fmt.Println("Init() is run.")
	fmt.Println("c:",c)
}

func (c *Controller) Test(){
	fmt.Println("Test() is run.")
}

//反射包练习
func main(){
	//初始化
	runController := &Controller{
		Action:"Run1",
		Method:"GET",
	}

	//Controller实现了ControllerInterface方法,因此它就实现了ControllerInterface接口
	var i ControllerInterface
	i = runController

	// 1.得到实际的值,通过value获取存储在里面的值,也可以去改变值，是i接口为指向main包下struct Controller的指针(runController)目前的所存储值
	value := reflect.ValueOf(i)
	fmt.Println("value:", value)

	// Elem返回值v包含的接口
	controllerValue := value.Elem()
	fmt.Println("controllerType(reflect.Value):",controllerValue)
	//获取存储在第一个字段里面的值
	fmt.Println("Action:", controllerValue.Field(0).String())

	valueMethod := value.MethodByName("Init")
	fmt.Println(valueMethod)

	// 2.得到结构体指针，是i接口为指向main包下struct Controller的指针类型。
	t := reflect.TypeOf(i)
	fmt.Println("type:",t)
	// Elem返回类型的元素类型。
	controllerType := t.Elem()
	tag := controllerType.Field(2).Tag //Field(第几个字段,index从0开始)
	fmt.Println("Tag:", tag)

	method, _ := t.MethodByName("Init")
	fmt.Println(method)

	// 有输入参数的方法调用
	// 构造参数
	args1 := []reflect.Value{reflect.ValueOf("Run2"),reflect.ValueOf("POST")}
	// 进行调用
	value.MethodByName("Init").Call(args1)

	// 无输入参数的方法调用
	// 构造参数
	args2 := make([]reflect.Value, 0)
	// 通过v进行调用
	value.MethodByName("Test").Call(args2)

}

