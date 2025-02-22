import { RectResize } from "@logicflow/extension";

import { h } from 'preact';

// 提供节点
class HttpNode extends RectResize.view {
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
        y: y - height / 2 + 10,
        width: 800,
        height: 800,
        viewBox: "0 0 1024 1024",
      }, [
        h('path', {
          fill: style.stroke,
          d: "M17.303 6.697a7.5 7.5 0 0 0-12.47 3.088A5.002 5.002 0 0 0 6.5 19.5h12a4 4 0 0 0 .99-7.876a7.474 7.474 0 0 0-2.187-4.927",
        }),
      ])
    ]);
  }
}

// 提供节点的属性
class HttpNodeModel extends RectResize.model {
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
  type: "HttpNode",
  view: HttpNode,
  model: HttpNodeModel
}