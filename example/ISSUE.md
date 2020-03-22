# ISSUE LIST

业务处理目标

## OMS

* 订单接入 -- 校验

  * 订单渠道来源

    确定区分来路渠道

  * 订单规格类型

    确定不同来源规格形式

  * 订单校验

    确定 '协同' 系统/动作  - CRM[用户], 财务[支付流水], 商品[价量/促销]
  
    确定 '协同' 输出

* 订单接入 -- 确认

  * CRM[用户信息/行为确认]

  * 财务[确订订单入账]

* 订单管理

  * 多路由订单管理
  
  * 订单查询、管理接口（UI/API）支持

* 订单状态流转控制
  
  * 确定发起方/接收方及形式

    相关系统: WMS,库存系统,售后,商品系统

## 库管中心

* 库管订单管理

    * 接入单
     
        确定最优发货路径策略
        
        如：ToB批次发送，ToC分仓发送
    
    * 回滚/取消
    
        拦截出库
        
        拒收货品翻仓
        
        退货申请后，发货仓原路径返回
        
        多实体库动态库存实时汇总与呈现

* 调拨单

    * 调拨下单
        发起商品调拨

    * 调拨确认

    * 取消调拨

* 库存管理

    管理 可用，现货，锁定，预占 库存