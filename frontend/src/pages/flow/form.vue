<template>
  <a-card size="small" :title="$store.menus.subTitle">
    <paas-form ref="form" v-model:model="model" :rules="vRules" add="/fwconfig" edit="/fwconfig/edit"
      label-width="80px">
      <form-input label="流程名称" name="name" v-model="model.name" :size="1"></form-input>
      <form-select label="流程类型" name="flowType" v-model="model.flowType" params="flowType" :size="1"></form-select>
      <form-select label="是否共享" name="isOpen" v-model="model.isOpen" params="isOpen" :size="1"></form-select>
      <form-select label="流程状态" name="flowState" v-model="model.flowState" params="flowState" :size="1"></form-select>
      <form-textarea label="流程描述" name="remark" v-model="model.remark"></form-textarea>
    </paas-form>
  </a-card>
</template>

<script>
export default {
  async mounted() {
    let id = this.$route.query.id;
    if (id) {
      let res = await this.$get(`/fwconfig/${id}`);
      this.model = res.data.data || {};
    }
  },
  data() {
    return {
      model: {
        groupType: "",
        groupName: "",
        parentId: "",
      },
      vRules: {
        groupType: [{ required: true, message: "请选择分组类型" }],
        groupName: [{ required: true, message: "请输入分组名称" }],
      },
    };
  },
};
</script>
