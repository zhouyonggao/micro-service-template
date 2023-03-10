// Code generated by entimport, DO NOT EDIT.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("order_no").Unique().Comment("订单号"),
		field.Int("user_id").Comment("下单用户ID"),
		field.Int("product_id").Comment("商品ID"), field.String("product_name").Comment("商品名称"), field.String("product_logo").Comment("商品logo图片地址"), field.Float("price").Comment("商品价格"), field.Float("denomination").Comment("商品面额"), field.Float("original_price").Comment("商品原价"), field.Float("unit_price").Comment("商品单价"), field.Int("num").Comment("购买数量"), field.Int("product_type").Comment("商品类型"), field.Int("status").Comment("订单状态"), field.Time("u_delete_time").Optional().Comment("用户删除时间"), field.Time("create_time").Comment("创建时间"), field.Time("update_time").Comment("更新时间")}
}
func (Order) Edges() []ent.Edge {
	return nil
}
func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "order_copy_tmp"}}
}
