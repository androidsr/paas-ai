import { CircleNode, CircleNodeModel } from "@logicflow/core";
import { h } from 'preact';

// 提供节点
class StartNode extends CircleNode {
  getIconShape () {
    const { model } = this.props
    const {
      x,
      y,
      width,
      height
    } = model
    const stroke = '#1890FF'
    return h(
      'svg',
      {
        x: x - width / 2,
        y: y - height / 2,
        width: 80,
        height: 80,
        viewBox: '0 0 1024 1024'
      },
      h(
        'path',
        {
          fill: stroke,
          d: 'M320 176a16 16 0 0 0-16 16v53l-111.68-67.44a10.78 10.78 0 0 0-16.32 9.31v138.26a10.78 10.78 0 0 0 16.32 9.31L304 267v53a16 16 0 0 0 32 0V192a16 16 0 0 0-16-16'
        }
      )
    )
  }
  getShape () {
    const { model } = this.props
    const { x, y, r } = model
    const { fill, stroke, strokeWidth } = model.getNodeStyle()
    return h(
      'g',
      {
      },
      [
        h(
          'circle',
          {
            cx: x,
            cy: y,
            r,
            fill,
            stroke,
            strokeWidth
          }
        ),
        this.getIconShape()
      ]
    )
  }
}

// 提供节点的属性
class StartNodeModel extends CircleNodeModel {
  // 自定义节点形状属性
  initNodeData (data) {
    data.text = {
      value: (data.text && data.text.value) || '',
      x: data.x,
      y: data.y + 35
    }
    super.initNodeData(data)
    this.r = 20
  }
  // 自定义节点样式属性
  getNodeStyle () {
    const style = super.getNodeStyle()
    return style
  }
  // 自定义锚点样式
  getAnchorStyle () {
    const style = super.getAnchorStyle();
    style.hover.r = 8;
    //style.hover.fill = "rgb(24, 125, 255)";
    //style.hover.stroke = "rgb(24, 125, 255)";
    return style;
  }
  // 自定义节点outline
  getOutlineStyle () {
    const style = super.getOutlineStyle();
    //style.stroke = '#88f'
    return style
  }
  getConnectedTargetRules () {
    const rules = super.getConnectedTargetRules()
    const notAsTarget = {
      message: '起始节点不能作为连线的终点',
      validate: () => false
    }
    rules.push(notAsTarget)
    return rules
  }
}

export default {
  type: "StartNode",
  view: StartNode,
  model: StartNodeModel
}