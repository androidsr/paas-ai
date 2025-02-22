import { RectResize } from "@logicflow/extension";

import { h } from 'preact';

// 提供节点
class VectorNode extends RectResize.view {
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
        y: y - height / 2 + 8,
        width: 40,
        height: 40,
        viewBox: "0 0 1024 1024",
      }, [
        h('path', {
          fill: style.stroke,
          d: "m479.66 268.7l-32-151.81C441.48 83.77 417.68 64 384 64H128c-16.8 0-31 4.69-42.1 13.94s-18.37 22.31-21.58 38.89l-32 151.87A16.7 16.7 0 0 0 32 272v112a64 64 0 0 0 64 64h320a64 64 0 0 0 64-64V272a16.7 16.7 0 0 0-.34-3.3m-384-145.4v-.28c3.55-18.43 13.81-27 32.29-27H384c18.61 0 28.87 8.55 32.27 26.91c0 .13.05.26.07.39l26.93 127.88a4 4 0 0 1-3.92 4.82H320a15.92 15.92 0 0 0-16 15.82a48 48 0 1 1-96 0A15.92 15.92 0 0 0 192 256H72.65a4 4 0 0 1-3.92-4.82Z",
        }),
      ])

    ]);
  }
}

// 提供节点的属性
class VectorNodeModel extends RectResize.model {
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
  type: "VectorNode",
  view: VectorNode,
  model: VectorNodeModel
}