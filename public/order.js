async function handleOrder() {
  const orderData = {
    customer_id: "123123",
    Items: [
      {"ID":"prod_RIBRYeGX7Vpawu","Quantity":2},
      {"ID":"prod_RIBRYeGX7Vpawu","Quantity":1},
      {"ID":"prod_RIBRYeGX7Vpawu","Quantity":5},
      {"ID":"prod_RIBRYeGX7Vpawu","Quantity":2}
    ]
  };

  try {
    const response = await fetch(`http://localhost:8282/api/customer/${orderData.customer_id}/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(orderData)
    });

    const result = await response.json();

    if (result.errno === 0) {
      // 跳转到success页面
      window.location.href = result.data.redirect_url;
    } else {
      alert('下单失败：' + result.message);
    }
  } catch (error) {
    console.error('下单出错：', error);
    alert('下单发生错误，请重试');
  }
} 