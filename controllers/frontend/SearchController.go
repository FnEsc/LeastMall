package frontend

import (
	"LeastMall/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
	"math"
	"reflect"
	"strconv"
)

type SearchController struct {
	BaseController
}

func (c SearchController) CreateProduct() {
	var product []models.Product
	models.DB.Find(&product)

	for i := 0; i < len(product); i++ {
		res, err := models.EsClient.Index().
			Index("product").
			Id(strconv.Itoa(product[i].Id)).
			BodyJson(product[i]).
			Do(context.Background())
		if err != nil {
			logs.Error("CreateProduct Failed", err)
		}
		fmt.Printf("CreateProduct %s\n", res.Result)
	}
	c.Ctx.WriteString("CreateProduct Success")
}

func (c *SearchController) UpdateProduct() {
	// 从数据库获取秀爱
	var product models.Product
	models.DB.Where("id=20").Find(&product)
	product.Title = "苹果电脑标题"
	product.SubTitle = "苹果电脑子标题"
	res, err := models.EsClient.Update().Index("product").Id("20").Doc(product).Do(context.Background())
	if err != nil {
		logs.Error("UpdateProduct Failed", err)
	}
	fmt.Printf("UpdateProduct %s\n", res.Result)
	c.Ctx.WriteString("UpdateProduct Success")
}

func (c *SearchController) DeleteProduct() {
	res, err := models.EsClient.Delete().Index("product").Id("20").Do(context.Background())
	if err != nil {
		logs.Error("DeleteProduct Failed", err)
	}
	fmt.Printf("DeleteProduct %s\n", res.Result)
	c.Ctx.WriteString("DeleteProduct Success")
}

func (c *SearchController) GetOne() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			c.Ctx.WriteString("GetOne Recovered")
		}
	}()
	result, _ := models.EsClient.Get().Index("product").Id("19").Do(context.Background())
	fmt.Println(result.Source)

	var product models.Product
	json.Unmarshal(result.Source, &product)
	c.Data["json"] = product
	c.ServeJSON()
}

// Query 查询多条数据
func (c *SearchController) Query() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			c.Ctx.WriteString("Query Recovered")
		}
	}()

	query := elastic.NewMatchQuery("Title", "旗舰")
	searchResult, err := models.EsClient.Search().Index("product").Query(query).Do(context.Background())
	if err != nil {
		panic(err)
	}
	var productList []models.Product
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		g := item.(models.Product)
		fmt.Printf("ID[%d], 标题[%v]\n", g.Id, g.Title)
		productList = append(productList, g)
	}
	c.Data["json"] = productList
	c.ServeJSON()
}

func (c *SearchController) FilterQuery() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			c.Ctx.WriteString("FilterQuery Recovered")
		}
	}()

	// 筛选
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("Title", "小米"))
	boolQ.Filter(elastic.NewRangeQuery("Id").Gt(19))
	boolQ.Filter(elastic.NewRangeQuery("Id").Lt(31))
	searchResult, err := models.EsClient.Search().Index("product").Query(boolQ).Do(context.Background())
	if err != nil {
		fmt.Println("FilterQuery Error", err)
	}
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		g := item.(models.Product)
		fmt.Printf("ID[%d], 标题[%v]\n", g.Id, g.Title)
	}
	c.Ctx.WriteString("FilterQuery Success")
}

// ProductList 分页查询
func (c *SearchController) ProductList() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			c.Ctx.WriteString("FilterQuery Recovered")
		}
	}()

	c.BaseInit()
	keyword := c.GetString("keyword")

	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	pageSize := 5
	query := elastic.NewMatchQuery("Title", keyword)
	searchResult, err := models.EsClient.Search().
		Index("product").
		Query(query).
		Sort("Price", true).
		Sort("Id", false).
		From((page - 1) * pageSize).
		Size(pageSize).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	// 查询符合条件的商品的总数
	searchResult2, _ := models.EsClient.Search().Index("product").Query(query).Do(context.Background())
	var productList []models.Product
	var product models.Product
	for _, item := range searchResult.Each(reflect.TypeOf(product)) {
		g := item.(models.Product)
		fmt.Printf("ID[%d], 标题[%v]\n", g.Id, g.Title)
		productList = append(productList, g)
	}
	c.Data["productList"] = productList
	c.Data["totalPages"] = math.Ceil(float64(len(searchResult2.Each(reflect.TypeOf(product)))) / float64(pageSize))
	c.Data["page"] = page
	c.Data["keyword"] = keyword
	c.TplName = "frontend/elasticsearch/list.html"

}
