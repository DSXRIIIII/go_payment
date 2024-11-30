package adapters

import (
	"context"
	domain "github.com/dsxriiiii/l3x_pay/order/domain/order"
	"github.com/dsxriiiii/l3x_pay/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TODO 申明全局变量有误
//var (
//	dbName   = viper.GetString("mongo.db-name")
//	collName = viper.GetString("mongo.coll-name")
//)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

type orderModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) collection() *mongo.Collection {
	return r.db.Database(viper.GetString("mongo.db-name")).Collection(viper.GetString("mongo.coll-name"))
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (r OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer r.logWithTag("create", err, order, created)
	write := r.marshalToModel(order)
	res, err := r.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return created, nil
}

func (r OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	defer r.logWithTag("get", err, nil, got)
	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	cond := bson.M{"_id": mongoID}
	logrus.Infof("cond:%s,read:%s", cond, &read)
	if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
		return
	}
	if read == nil {
		return nil, domain.NotFoundError{
			OrderID: id,
		}
	}
	got = r.unmarshal(read)
	return got, nil
}

func (r OrderRepositoryMongo) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) (err error) {
	defer r.logWithTag("update", err, order, nil)
	if order == nil {
		panic("got nil order")
	}
	session, err := r.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)
	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = session.CommitTransaction(ctx)
		} else {
			_ = session.AbortTransaction(ctx)
		}
	}()

	/** start transaction **/
	oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, order)
	if err != nil {
		return
	}
	logrus.Infof("update||oldOrder=%+v||updated=%+v", oldOrder, updated)
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	res, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	r.logWithTag("finish_update", err, order, res)
	return
}

func (r *OrderRepositoryMongo) logWithTag(tag string, err error, input *domain.Order, result interface{}) {
	l := logrus.WithFields(logrus.Fields{
		"tag":            "order_repository_mongo",
		"input_order":    input,
		"performed_time": time.Now().Unix(),
		"err":            err,
		"result":         result,
	})
	if err != nil {
		l.Infof("%s_fail", tag)
	} else {
		l.Infof("%s_success", tag)
	}
}

func (r *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
	return &orderModel{
		MongoID:     primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
	return &domain.Order{
		ID:          m.MongoID.Hex(),
		CustomerID:  m.CustomerID,
		Status:      m.Status,
		PaymentLink: m.PaymentLink,
		Items:       m.Items,
	}
}
