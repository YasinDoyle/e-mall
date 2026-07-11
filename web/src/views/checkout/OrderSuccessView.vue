<template>
  <div style="text-align: center; padding: 60px">
    <el-result
      icon="success"
      title="支付成功！"
      sub-title="感谢您的购买，商家将尽快发货"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/user/orders')"
          >查看订单</el-button
        >
        <el-button @click="$router.push('/')">继续购物</el-button>
      </template>
    </el-result>
  </div>
</template>
