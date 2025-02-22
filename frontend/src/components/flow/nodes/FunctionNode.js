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

    return h('g', {}, [
      h('rect', {
        ...style,
        x: x - width / 2,
        y: y - height / 2,
        width,
        height,
        rx: radius,
        ry: radius,
      }),
      h('svg', {
        x: x - width / 2 + 3,
        y: y - height / 2 + 5,
        width: 40,
        height: 50,
        viewBox: "0 0 1024 1024",
      }, [
        h('path', {
          fill: style.stroke,
          d: "M472.9 71c-4.5-4.5-10.7-7-17-7H56.2c-6.4 0-12.5 2.5-17 7S32 81.7 32 88v272c0 6.4 2.7 12.5 7.2 17s10.5 7 16.8 7h400c6.4 0 12.3-2.5 16.8-7s7.2-10.7 7.2-17V88c.1-6.3-2.6-12.5-7.1-17zM256 360c-4.4 0-8-3.6-8-8s3.6-8 8-8 8 3.6 8 8-3.6 8-8 8zm203.9-40H52.2c-2.2 0-4-1.8-4-4V92c0-6.6 5.4-12 12-12h391.7c6.6 0 12 5.4 12 12v224c0 2.2-1.8 4-4 4z",
        }),
      ])

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
       this.height = 36;
    this.radius = 2;
  }
  getNodeStyle () {
    const style = super.getNodeStyle();
    return style;
  }
}

export default {
  type: "FunctionNode",
  view: FunctionNode,
  model: FunctionNodeModel
}