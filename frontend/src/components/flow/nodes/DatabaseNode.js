import { RectResize } from "@logicflow/extension";

import { h } from 'preact';

// 提供节点
class DatabaseNode extends RectResize.view {
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
          d: "M1.25 9A6.75 6.75 0 0 1 8 2.25h8A6.75 6.75 0 0 1 22.75 9v6A6.75 6.75 0 0 1 16 21.75H8A6.75 6.75 0 0 1 1.25 15zm16.08-.096c-.25-.14-.588-.154-.973-.154H14.25v2.5h2.128c.377 0 .706-.017.952-.154c.158-.089.42-.302.42-1.096s-.262-1.007-.42-1.096M18.582 12c.423-.45.668-1.112.668-2c0-1.206-.452-1.993-1.187-2.404c-.62-.348-1.325-.346-1.67-.346H13.5a.75.75 0 0 0-.75.75v8c0 .414.336.75.75.75h2.894c.344 0 1.049.002 1.669-.346c.735-.411 1.187-1.198 1.187-2.404c0-.888-.245-1.549-.668-2m-2.203.75H14.25v2.5h2.107c.385 0 .723-.014.973-.154c.158-.089.42-.302.42-1.096s-.262-1.007-.42-1.095c-.245-.138-.575-.154-.951-.155m-8.022-4c.385 0 .723.014.973.154c.158.089.42.302.42 1.096v4c0 .794-.262 1.007-.42 1.096c-.25.14-.588.154-.973.154H6.25v-6.5zM11.25 10c0-1.206-.452-1.993-1.187-2.404c-.62-.348-1.325-.346-1.67-.346H5.5a.75.75 0 0 0-.75.75v8c0 .414.336.75.75.75h2.894c.344 0 1.049.002 1.669-.346c.735-.411 1.187-1.198 1.187-2.404z",
        }),
      ])

    ]);
  }
}

// 提供节点的属性
class DatabaseNodeModel extends RectResize.model {
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
  type: "DatabaseNode",
  view: DatabaseNode,
  model: DatabaseNodeModel
}