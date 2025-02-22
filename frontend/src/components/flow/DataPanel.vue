<template>
  <div class="data-panel buttons">
    <a-button-group>
      <a-button @click="() => { this.lf.zoom(true) }"><PlusSquareOutlined /></a-button>
      <a-button @click="() => { this.lf.zoom(false) }"><MinusSquareOutlined /></a-button>
      <a-button @click="() => { this.lf.resetZoom() }"><RedoOutlined /></a-button>
      <a-button @click="() => { this.lf.undo() }"><SwapLeftOutlined /></a-button>
      <a-button @click="() => { this.lf.redo() }"><SwapRightOutlined /></a-button>
      <a-button @click="downloadImage"><AreaChartOutlined /></a-button>
      <a-button @click="show"><BoxPlotOutlined /></a-button>
      <a-button @click="save" type="primary"><CheckOutlined /></a-button>
      <a-button @click="() => { this.$router.go(-1) }" type="primary" danger><CloseOutlined /></a-button>
    </a-button-group>
    <a-modal v-model:open="opend" :closable="true" :maskClosable="true" :draggable="true" title="查看数据">
      <a-textarea v-model:value="content" :rows="25" style="width: 100%;"></a-textarea>
      <!-- 自定义 footer -->
      <template #footer>
        <a-button @click="opend = false">取消</a-button>
        <a-button type="primary" @click="applyChanges">确认</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts">
import { message } from 'ant-design-vue';

export default {
  name: "DataPanel",
  mounted() {
    this.id = this.$route.query.id;
    this.$get("/fwconfig/" + this.id).then(res => {
      this.content = res.data.data.content;
      if (!this.content) {
        return;
      }
      let jsonData = JSON.parse(this.content);
      for (let i = 0; i < jsonData["nodes"].length; i++) {
        let item = jsonData["nodes"][i];
        delete item.properties.width;
        delete item.properties.height;
      }
      (this.$props.lf).render(jsonData);
    });
  },
  data() {
    return {
      opend: false,
      content: "",
      id: "",
    };
  },
  props: {
    lf: Object,

  },

  methods: {
    downloadImage() {
      (this.$props.lf).getSnapshot();
    },
    save() {
      const data = JSON.stringify((this.$props.lf).getGraphRawData());
      this.$put("/fwconfig/edit", { id: this.id, content: data }).then(res => {
        if (res.data.code == this.$success) {
          message.success(res.data.msg);
        }
      });
    },
    show() {
      this.content = JSON.stringify((this.$props.lf).getGraphRawData(), null, '\t');
      this.opend = true;
    },
    applyChanges() {
      try {
        const parsedData = JSON.parse(this.content);
        this.$props.lf.render(parsedData);
        message.success("数据已应用到流程图");
        this.opend = false;
      } catch (error) {
        console.error(error);
        message.error("JSON 格式有误，请检查！");
      }
    },
  },
};
</script>
<style>
.data-panel {
  margin: 10px;
  position: absolute;
  top: 10;
  width: auto;
  background-color: rgb(227, 242, 248);
  box-shadow: 0 0 10px 1px rgb(228, 224, 219);
  z-index: 101;
}
</style>