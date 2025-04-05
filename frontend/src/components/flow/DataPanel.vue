<template>
  <div class="data-panel">
    <a-button-group>
      <a-button @click="() => { this.lf.resetZoom() }">适合</a-button>
      <a-button @click="() => { this.lf.zoom(true) }">放大</a-button>
      <a-button @click="() => { this.lf.zoom(false) }">缩小</a-button>
      <a-button @click="() => { this.lf.undo() }">上一步</a-button>
      <a-button @click="() => { this.lf.redo() }">下一步</a-button>
      <a-button @click="downloadImage">导出图片</a-button>
      <a-button @click="show">查看数据</a-button>
      <a-button @click="save">保存</a-button>
      <a-button @click="() => { this.$router.go(-1) }">关闭</a-button>
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

<script>
import { message } from 'ant-design-vue';
import { Edit, Get } from '../../../wailsjs/go/biz/FwConfigBiz.js';

export default {
  name: "DataPanel",
  mounted() {
    this.id = this.$route.query.id;
    Get(this.id).then(res => {
      this.content = res.data.content;
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
      Edit({ id: this.id, content: data }).then(res => {
        if (res.code == 200) {
          this.$message.success(res.msg);
          this.back();
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
  position: absolute;
  top: 10px;
  right: 10px;
  width: auto;
  box-shadow: 0 0 10px 1px rgb(228, 224, 219);
  z-index: 101;
  display: flex;
}
</style>