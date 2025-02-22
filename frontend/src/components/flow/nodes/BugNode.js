import { RectResize } from "@logicflow/extension";

import { h } from 'preact';

// 提供节点
class BugNode extends RectResize.view {
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
        width: 900,
        height: 900,
        viewBox: "0 0 1024 1024",
      }, [
        h('path', {
          fill: style.stroke,
          d: "M12 21q-1.6 0-2.863-.825T7.076 18l-1.525.875q-.35.2-.75.075t-.6-.475t-.1-.75t.45-.6l1.725-1q-.075-.275-.125-.563T6.05 15H4q-.425 0-.712-.288T3 14t.288-.712T4 13h2.05q.05-.3.1-.587t.125-.563l-1.725-1q-.35-.2-.45-.6t.1-.75t.613-.462t.762.087L7.05 10q.2-.35.463-.687T8.05 8.7Q8 8.525 8 8.35V8q0-.6.175-1.15t.475-1.025l-.95-.95q-.275-.275-.288-.7T7.7 3.45q.275-.3.688-.287t.712.287l1.05 1q.425-.225.888-.337T12 4t.975.125t.9.35L14.9 3.45q.3-.3.7-.288t.7.313q.275.3.288.7t-.288.7l-.95.95q.3.475.463 1.025T15.975 8v.338q0 .162-.05.337q.275.275.538.625t.462.7l1.525-.875q.35-.2.75-.088t.6.463t.088.763t-.463.612l-1.725.975q.075.275.138.563t.112.587H20q.425 0 .713.288T21 14t-.288.713T20 15h-2.05q-.05.3-.1.588t-.125.562l1.725 1q.35.2.45.613t-.1.762t-.6.45t-.75-.1L16.925 18q-.8 1.35-2.063 2.175T12 21M10.1 7.35q.425-.175.913-.262T12 7t.963.075t.887.25q-.2-.575-.7-.95T12 6t-1.175.388t-.725.962M12 19q1.825 0 2.913-1.525T16 14q0-1.75-1.012-3.375T12 9q-1.95 0-2.975 1.613T8 14q0 1.95 1.088 3.475T12 19m0-2q-.425 0-.712-.288T11 16v-4q0-.425.288-.712T12 11t.713.288T13 12v4q0 .425-.288.713T12 17",
        }),
      ])

    ]);
  }
}

// 提供节点的属性
class BugNodeModel extends RectResize.model {
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
  type: "BugNode",
  view: BugNode,
  model: BugNodeModel
}