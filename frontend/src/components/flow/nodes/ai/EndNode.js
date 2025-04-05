import { RectResize } from "@logicflow/extension";
import { h } from 'preact';

// 提供节点
class EndNode extends RectResize.view {
  getResizeShape() {
    const { model } = this.props;
    const { x, y, width, height, radius, properties } = model;
    const style = model.getNodeStyle();

    const titleHeight = 26; // 标题栏高度

    return h('g', { filter: "drop-shadow(1px 1px 2px rgba(13, 15, 39, 0.1))" }, [
      // 背景矩形
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
        y: y - height / 2,
        width,
        height: titleHeight,
        fill: "#F44336",  // 红色背景
        stroke: "#F44336",
        strokeWidth: 1,
      }),

      // 标题文本
      h('text', {
        x: x - width / 2 + 10,
        y: y - height / 2 + titleHeight / 2 + 4,
        textAnchor: "start",
        fontSize: "14px",
        fontWeight: "bold",
        fill: "#ffffff",  // 白色标题文字
      }, properties.title),
    ]);
  }
}

// 提供节点的属性
class EndNodeModel extends RectResize.model {
  constructor(data, graphModel) {
    super(data, graphModel);
  }

  initNodeData(data) {
    super.initNodeData(data);
    this.width = 200;
    this.height = 90;  // 适应标题栏
    //this.radius = 5;
    this.text.draggable = false;
    this.text.editable = false;
    this.properties = {
      ...this.properties,
      title: "结束"
    };
  }

  getTextStyle() {
    const style = super.getTextStyle();
    style.fontSize = 15;
    style.fontWeight = "bold";
    style.color = "#F44336";
    return style;
  }

  getNodeStyle() {
    return {
      fill: "#fff",
      stroke: "#e0e0e0",
      strokeWidth: 1,
    };
  }

  getConnectedSourceRules() {
    const rules = super.getConnectedSourceRules();
    rules.push({
      message: '终止节点不能作为连线的起点',
      validate: () => false
    });
    return rules;
  }
}

export default {
  type: "EndNode",
  view: EndNode,
  model: EndNodeModel
};
