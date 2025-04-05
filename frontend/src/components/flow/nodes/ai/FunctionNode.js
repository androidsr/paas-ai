import { RectResize } from "@logicflow/extension";
import { h } from 'preact';

// 提供节点
class FunctionNode extends RectResize.view {
  getResizeShape() {
    // 获取XxxNodeModel中定义的形状属性
    const { model } = this.props;
    const { x, y, width, height, radius, properties } = model;
    // 获取XxxNodeModel中定义的样式属性
    const style = model.getNodeStyle();

    // 标题栏高度
    const titleHeight = 26;

    return h('g', { filter: "drop-shadow(1px 1px 2px rgba(13, 15, 39, 0.1))" }, [
      // 外层矩形（整体背景）
      h('rect', {
        ...style,
        x: x - width / 2,
        y: y - height / 2,
        width,
        height,
        rx: radius,
        ry: radius,
      }),

      // 标题栏
      h('rect', {
        x: x - width / 2,
        y: y - height / 2, // 让标题栏浮动，避免挡住圆角
        width,
        height: titleHeight,
        fill: "#82b366", // 天蓝色背景
        stroke: "#82b366",
        strokeWidth: 1,
      }),

      // 标题文本（靠左显示，文本在标题栏上下居中）
      h('text', {
        x: x - width / 2 + 10, // 保证文本距离左边缘一定距离
        y: y - height / 2 + titleHeight / 2 + 4, // 垂直居中
        textAnchor: "start",
        fontSize: "14px",
        fontWeight: "bold",
        fill: "#ffffff", // 略浅的文字
      }, properties.title),
    ]);
  }
}

// 提供节点的属性
class FunctionNodeModel extends RectResize.model {
  constructor(data, graphModel) {
    super(data, graphModel);
  }

  initNodeData (data) {
    super.initNodeData(data);
    this.width = 200;
    this.height = 90;  // 调整高度，以便显示完整的标题栏
    //this.radius = 2;
    this.text.draggable = false; // 不允许文本被拖动
    this.text.editable = false; // 不允许文本被编辑
    this.properties = {
      ...this.properties,
      title: "工具"
    };
  }

  getNodeStyle () {
    return {
      fill: "#fff",
      stroke: "#e0e0e0", // 更浅的边框
      strokeWidth: 1,
    };
  }
}

export default {
  type: "FunctionNode",
  view: FunctionNode,
  model: FunctionNodeModel
};
