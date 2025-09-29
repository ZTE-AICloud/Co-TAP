/* Started by AICoder, pid:3362bx3734pf32d1424e0956e03fbe7c4dc4c15f */
import { Component } from '@angular/core';
import { applyPatch } from 'fast-json-patch';
import { HttpAgent, Message } from 'hai-client';
import { AgentSubscriber } from 'hai-client/dist/client/agent/subscriber';
import { PlxDialogue, PlxMessageRoleType } from 'paletx/ai/plx-chat';
import { PlxContentType } from 'paletx/ai/plx-content';
import { HeaderAction } from 'paletx/ai/plx-frame-native/plx-frame-native.model';
import { RestAPI } from '../../shared/util/rest-api';

@Component({
  selector: 'task',
  templateUrl: 'task.component.html',
  styleUrls: ['./task.component.less']
})
export class TaskComponent {
  title = '张江高科3月20日上午出********';
  logo = "assets/img/logo.png";
  showTab = false;
  activeMode: 'chat' | 'workflow' = 'chat';
  tabActiveId = 'tabId1';
  tabItems = [
    {
      id: 'tabId1',
      title: '首页'
    },
    {
      id: 'tabId2',
      title: '故障诊断'
    }
  ];
  headerActions: HeaderAction[] = [
    {
      icon: 'plx-ico-sens-box-f-20',
      name: 'Sense Box',
      isWarning: true
    },
    {
      icon: 'plx-ico-detail-16',
      name: '智能体库',
    },
  ];
  state = { workflowItems: [] };
  // 对话
  dialogue: PlxDialogue;
  senderConfig: any;
  // 任务
  tasks = [];
  // 存储toolCall数据，用于去重
  toolCallMaps = new Map();
  activeToolStep: any;

  ngOnInit() {
    this.setDialogue();
  }

  onTabChange(event) {
    this.tabActiveId = event.nextId;
  }

  public bubbleContentAction(event: any): void {
    if (event?.data.name === '开始执行') {
      this.addUserMessageToDialogue(event.data.name);

      setTimeout(() => {
        this.showTab = true;
        this.activeMode = 'workflow';
        this.tabActiveId = 'tabId2';
        this.getTaskData();
      }, 500);
    }
  }

  /* Started by AICoder, pid:2950adb8ac604e914bde09e950aabc113e482b91 */
  private addUserMessageToDialogue(content: string): void {
    const messages = this.dialogue.messages || [];

    const contents = [{
      type: PlxContentType.TEXT,
      content,
    }];

    messages.push({
      id: new Date().getTime().toString(),
      role: PlxMessageRoleType.USER,
      name: 'admin',
      contents,
      timestamp: new Date().getTime(),
    });

    this.dialogue.messages = messages;
  }
  /* Ended by AICoder, pid:2950adb8ac604e914bde09e950aabc113e482b91 */

  private async getTaskData(): Promise<void> {
    this.tasks = [];
    this.state.workflowItems = [];
    const requestParam = {
      url: RestAPI.taskUrl,
      initialState: this.state,
    };
    const agent = new HttpAgent(requestParam);

    let currentTask: any;
    let currentStep: any;
    const subscriber: AgentSubscriber = {
      onRunStartedEvent: (eventData) => {
        console.log("RunStarted");
        const { event } = eventData;
        currentTask = this.taskHandle(this.tasks, event['messageId']);
      },
      onStepStartedEvent: (eventData) => {
        console.log("StepStarted");
        this.toolCallMaps.clear();
        const { event } = eventData;
        currentStep = this.stepHandle(currentTask, currentStep, event['stepName']);
      },
      onBusinessDataStartEvent: () => {
        console.log("BusinessDataStart");
      },
      onBusinessDataContentEvent: (eventData) => {
        console.log("BusinessDataContent:", eventData);
        const { textMessageBuffer } = eventData;
        this.handleTaskTextMessageBuffer(textMessageBuffer, currentStep);
      },
      onBusinessDataEndEvent: (eventData) => {
        console.log("BusinessDataEnd:", eventData);
        const { textMessageBuffer } = eventData;
        this.handleTaskTextMessageBuffer(textMessageBuffer, currentStep);
      },
      onAgentCollaborativeMessageStartEvent: (eventData) => {
        console.log("AgentCollaborativeMessageStart", eventData);
        const workflowItems = this.state.workflowItems;
        const toAgent = eventData.event['to'];
        if (toAgent === 'planerAgent') {
          workflowItems.push({avatar: 'assets/img/planner agent@4x.png'});
        } else if (toAgent === 'actorAgent') {
          workflowItems.push({avatar: 'assets/img/actor agent@4x.png'});
        }
      },
      onAgentCollaborativeMessageContentEvent: () => {
        console.log("AgentCollaborativeMessageContent");
      },
      onAgentCollaborativeMessageEndEvent: (eventData) => {
        console.log("AgentCollaborativeMessageEnd:", eventData);
        const { textMessageBuffer } = eventData;
        const item = JSON.parse(textMessageBuffer);
        const { task: { id, name } } = item;
        const workflowItems = this.state.workflowItems;
        const workflowItem = workflowItems.find(item => item.id === id);
        if (workflowItem) {
          workflowItem.name = name;
        } else {
          const lastItem = workflowItems[workflowItems.length - 1];
          if (lastItem && !lastItem.id) {
            lastItem.id = id;
            lastItem.name = name;
          } else {
            workflowItems.push({ id, name });
          }
        }
      },
      onToolCallStartEvent: () => {
        console.log("ToolCallStart");
      },
      onToolCallArgsEvent: (eventData) => {
        console.log("ToolCallArgs:", eventData);
        const deltaObj = this.parseJsonString(eventData.event.delta);
        deltaObj.length && this.ToolCallArgHandle(deltaObj[0].id);
      },
      onToolCallEndEvent: () => {
        console.log("ToolCallEnd");
        this.unActiveToolStep();
      },
      onStateDeltaEvent: (eventData) => {
        console.log("StateDelta:", eventData);
        const { event: { delta } } = eventData;
        try {
          // 同步状态数据
          const result = applyPatch(this.state, delta, true, false);
          this.state = result.newDocument;
          this.state.workflowItems.forEach(item => {
            if (item.status === 'running') {
              item.status = 'active';
            } else if (item.status === 'done') {
              item.status = 'success';
            }
          });
        } catch (e) {
          console.warn('applyPatch error', e);
        }
      },
      onStepFinishedEvent: () => {
        console.log("StepFinished");
        currentStep = null;
      },
      onRunFinishedEvent: () => {
        console.log("RunFinished");
        currentTask = null;
      },
    }

    const subscription = agent.subscribe(subscriber);
    const result = await agent.runAgent();
    console.log(result);
  }

  private unActiveToolStep() {
    this.activeToolStep && (this.activeToolStep.active = false);
  }

  private ToolCallArgHandle(toolStepId: string) {
    for (const [key, value] of this.toolCallMaps) {
      for (let toolStep of value.tools) {
        if (toolStepId === toolStep.id) {
          toolStep.active = true;
          toolStep.selected = true;
          this.activeToolStep = toolStep;
        }
      }
    }
  }

  /* Started by AICoder, pid:87209q9c13z46ff1480908bce07fa21a062001c1 */
  private handleTaskTextMessageBuffer(textMessageBuffer: string, currentStep: any): void {
    if (textMessageBuffer) {
      const dataList = this.parseJsonString(textMessageBuffer);
      const lastItem = dataList[dataList.length - 1];
      const lastOutput = lastItem?.output;
      if (lastOutput && lastOutput.type === 'think') {
        lastOutput.type = 'thinking';
      }

      if (lastItem?.blockId) {
        this.blockHandle(currentStep, lastOutput, lastItem?.blockId);
      } else {
        this.dataContentHandle(currentStep, lastOutput);
      }
    }
  }
  /* Ended by AICoder, pid:87209q9c13z46ff1480908bce07fa21a062001c1 */

  private taskHandle(tasks: any[], messageId: string): any {
    const currentTask = {
      messageId: messageId,
      steps: []
    };
    tasks.push(currentTask);
    return currentTask;
  }

  /* Started by AICoder, pid:s38d4pd275n3d9514bb80b09d000cc1340b8cd9d */
  private stepHandle(currentTask: any, currentStep: any, stepName: string): any {
    // 如果已有有效步骤直接返回
    if (currentStep) {
      return currentStep;
    }

    // 创建新步骤对象
    const newStep = {
      name: stepName,
      showAIGenerator: true,
      info: []
    };

    // 自动关联到任务（如果存在）
    if (currentTask && Array.isArray(currentTask.steps)) {
      currentTask.steps.push(newStep);
    }

    return newStep;
  }
  /* Ended by AICoder, pid:s38d4pd275n3d9514bb80b09d000cc1340b8cd9d */

  private dataContentHandle(currentStep: any, output: any): void {
    if (!currentStep) {
      return;
    }

    const contentType = output.type;
    if (contentType === 'thinking') {
      this.thinkingHanlde(currentStep, output);
    } else if (contentType === 'text') {
      this.textHandle(currentStep, output);
    } else if (contentType === 'toolCallShow') {
      this.toolHandle(currentStep, output);
    } else {
      if (contentType === 'fold' && output.data.position === 'end') {
        return;
      }
      currentStep.info.push({ type: contentType, data: output.data });
    }
  }

  private blockHandle(currentStep: any, output: any, blockId: string) {
    const lastItem = currentStep.info[currentStep.info.length - 1];
    let currentBlock = { type: 'block', blockId, info: [] };
    if (lastItem?.type === 'block' && blockId === lastItem.blockId) {
      currentBlock = lastItem;
    } else {
      currentStep.info.push(currentBlock);
    }
    this.dataContentHandle(currentBlock, output);
  }

  private thinkingHanlde(currentStep: any, output: any): void {
    // 处理思考中文本叠加
    const lastItem = currentStep.info[currentStep.info.length - 1];
    if (lastItem?.type === 'thinking') {
      lastItem.data.content += output.data.content;
      lastItem.data.collapse = output.data.collapse;
      lastItem.data.thinking = output.data.thinking;
    } else {
      currentStep.info.push({ type: 'thinking', data: output.data });
    }
  }

  private toolHandle(currentStep: any, output: any): void {
    if (this.toolCallMaps.get(output.id)) {
      const toolCall = this.toolCallMaps.get(output.id);
      toolCall.title = output.data.title;
      toolCall.icon = output.data.icon;
      output.data.tools.forEach(newTool => {
        const index = toolCall.tools?.findIndex(oldTool => oldTool.id === newTool.id);
        if (index > -1) {
          toolCall.tools[index] = newTool;
        } else {
          toolCall.tools.push(newTool);
        }
      });
    } else {
      const toolCall = {
        type: 'toolCallShow',
        id: output.id,
        title: output.data.title,
        icon: output.data.icon,
        tools: output.data.tools
      };
      this.toolCallMaps.set(output.id, toolCall);
      currentStep.info.push(toolCall);
    }
  }

  private textHandle(currentStep: any, output: any): void {
    // 处理连续文本叠加
    const lastItem = currentStep.info[currentStep.info.length - 1];
    if (lastItem?.type === 'text') {
      lastItem.data += output.content;
    } else {
      currentStep.info.push({ type: 'text', data: output.content });
    }
  }

  private setDialogue(): void {
    this.dialogue = new PlxDialogue();
    this.senderConfig = {
      value: '',
      disabled: false,
      submitDisabled: false,
      showStop: false,
      showOnlineSearchBtn: false,
    };
  }

  public async onSend(dialogue: PlxDialogue): Promise<void> {
    const dialogueMessages = dialogue.messages;
    const initialMessages = dialogueMessages.map(msg => {
      const { id, name, role, timestamp, contents } = msg;
      const content = contents.map(content => content.content).join('');
      return { id, name, role, timestamp, content };
    }) as Message[];

    const requestParam = {
      url: RestAPI.chatUrl,
      initialMessages,
    };

    const agent = new HttpAgent(requestParam);

    const subscriber: AgentSubscriber = {
      onBusinessDataStartEvent: () => {
        console.log("BusinessDataStart");
      },
      onBusinessDataContentEvent: (eventData) => {
        this.updateDialogueMessages(eventData, dialogueMessages);
      },
      onBusinessDataEndEvent: (eventData) => {
        this.updateDialogueMessages(eventData, dialogueMessages);
      },
    }
    const subscription = agent.subscribe(subscriber);

    const result = await agent.runAgent();
    console.log(result);
  }

  /* Started by AICoder, pid:69196jf783q7fca145cf0a13a05a8240dd9260c6 */
  private updateDialogueMessages(eventData: any, dialogueMessages: any[]): void {
    const { agent, event, messages } = eventData;
    const newAssistantMessage = messages[messages.length - 1];
    const { id, role, content } = newAssistantMessage || {};
    const dataList = this.parseJsonString(content);

    const newList = this.handleThinking(dataList);
    const contents = newList.map(d => {
      const { type, data, content } = d.output || {};
      return {
        type,
        data,
        content,
      };
    });

    let newDialogueMessage = dialogueMessages.find(
      message => message.id === (id + agent.agentId),
    );

    if (newDialogueMessage) {
      newDialogueMessage.contents = contents;
    } else {
      newDialogueMessage = {
        id: id + agent.agentId,
        role,
        avatar: 'assets/img/avatar.png',
        name: '智能体名称',
        contents,
        timestamp: event.timestamp,
      };
      dialogueMessages.push(newDialogueMessage);
    }

    this.dialogue = {
      ...this.dialogue,
      messages: dialogueMessages,
    };
  }
  /* Ended by AICoder, pid:69196jf783q7fca145cf0a13a05a8240dd9260c6 */

  /* Started by AICoder, pid:f19e8g3ac7ma670147300ba700192347589954ff */
  private handleThinking(dataList: any[]): any[] {
    let thinkingData: any;
    const newList = [];
    const THINKING_TYPE = 'thinking';

    for (const item of dataList) {
      if (!item?.output) {
        continue;
      }

      const { output } = item;

      if (output.type === THINKING_TYPE) {
        const thinkingOutput = output;

        if (thinkingData) {
          const existingData = thinkingData.output.data;
          const newData = thinkingOutput.data;

          // 合并内容
          existingData.content += newData.content;

          // 选择性更新字段（避免覆盖已有值）
          if (newData.collapse !== undefined) {
            existingData.collapse = newData.collapse;
          }
          if (newData.thinking !== undefined) {
            existingData.thinking = newData.thinking;
          }
        } else {
          // 深拷贝以避免修改原始数据
          thinkingData = {
            output: {
              ...thinkingOutput,
              data: {
                ...thinkingOutput.data
              }
            }
          };
        }
      } else {
        newList.push(item);
      }
    }

    if (thinkingData) {
      newList.unshift(thinkingData);
    }

    return newList;
  }
  /* Ended by AICoder, pid:f19e8g3ac7ma670147300ba700192347589954ff */

  private parseJsonString(jsonString: string): any[] {
    if (!jsonString || typeof jsonString !== 'string') {
      throw new Error('输入必须是有效的字符串');
    }

    try {
      // 处理原始字符串为合法JSON数组
      const legalJsonArrayStr = `[${jsonString.replace(/}{/g, '},{')}]`;

      // 解析并返回数组
      return JSON.parse(legalJsonArrayStr);
    } catch (error) {
      throw new Error(`JSON解析失败: ${(error as Error).message}`);
    }
  }
}
/* Ended by AICoder, pid:3362bx3734pf32d1424e0956e03fbe7c4dc4c15f */