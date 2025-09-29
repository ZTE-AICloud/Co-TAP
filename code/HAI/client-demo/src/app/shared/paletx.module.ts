import { NgModule } from '@angular/core';
import { PlxAiGeneratorModule } from 'paletx/ai/plx-ai-generator';
import { PlxBubbleModule } from 'paletx/ai/plx-bubble';
import { PlxChatModule } from 'paletx/ai/plx-chat';
import { PlxContentModule } from 'paletx/ai/plx-content';
import { PlxFrameNativeModule } from 'paletx/ai/plx-frame-native';
import { PlxSenderModule } from 'paletx/ai/plx-sender';
import { PlxTaskModule } from 'paletx/ai/plx-task';
import { PlxThinkingModule } from 'paletx/ai/plx-thinking';
import { PlxBreadcrumbModule } from 'paletx/plx-breadcrumb';
import { PlxButtonsModule } from 'paletx/plx-buttons';
import { PlxChartsModule } from 'paletx/plx-charts';
import { PlxDragHandleModule } from 'paletx/plx-drag-handle';
import { PlxDropdownModule } from 'paletx/plx-dropdown';
import { PlxFormModule } from 'paletx/plx-form';
import { PlxFormXModule } from 'paletx/plx-form-x';
import { PlxInfoModule } from 'paletx/plx-info';
import { PlxLoadingModule } from 'paletx/plx-loading';
import { PlxMessageModule } from 'paletx/plx-message';
import { PlxModalModule } from 'paletx/plx-modal';
import { PlxPopoverModule } from 'paletx/plx-popover';
import { PlxSelectModule } from 'paletx/plx-select';
import { PlxSystemPromptModule } from 'paletx/plx-systemprompt';
import { PlxTableModule } from 'paletx/plx-table';
import { PlxTabsetModule } from 'paletx/plx-tabSet';
import { PlxTextInputModule } from 'paletx/plx-text-input';
import { PlxToggleModule } from 'paletx/plx-toggle';
import { PlxToolbarModule } from 'paletx/plx-toolbar';
import { PlxTooltipModule } from 'paletx/plx-tooltip';
import { PlxViewModule } from 'paletx/plx-view';
import { PlxWorkflowModule } from 'paletx/ai/plx-workflow';

@NgModule({
  imports: [
    PlxBreadcrumbModule.forRoot(),
    PlxDropdownModule.forRoot(),
    PlxFormModule.forRoot(),
    PlxFormXModule.forRoot(),
    PlxInfoModule.forRoot(),
    PlxLoadingModule.forRoot(),
    PlxMessageModule.forRoot(),
    PlxModalModule.forRoot(),
    PlxPopoverModule.forRoot(),
    PlxSelectModule.forRoot(),
    PlxSystemPromptModule.forRoot(),
    PlxTextInputModule.forRoot(),
    PlxToggleModule.forRoot(),
    PlxTooltipModule.forRoot(),
    PlxToolbarModule.forRoot(),
    PlxViewModule.forRoot(),
    PlxButtonsModule.forRoot(),
    PlxDragHandleModule.forRoot(),
    PlxTableModule.forRoot(),
    PlxChartsModule.forRoot(),
    PlxTabsetModule.forRoot(),
    // AI组件
    PlxAiGeneratorModule.forRoot(),
    PlxBubbleModule.forRoot(),
    PlxChatModule.forRoot(),
    PlxContentModule.forRoot(),
    PlxFrameNativeModule.forRoot(),
    PlxSenderModule.forRoot(),
    PlxTaskModule.forRoot(),
    PlxThinkingModule.forRoot(),
    PlxWorkflowModule.forRoot(),
  ],
  exports: [
    PlxBreadcrumbModule,
    PlxDropdownModule,
    PlxFormModule,
    PlxFormXModule,
    PlxInfoModule,
    PlxLoadingModule,
    PlxMessageModule,
    PlxModalModule,
    PlxPopoverModule,
    PlxSelectModule,
    PlxSystemPromptModule,
    PlxTextInputModule,
    PlxToggleModule,
    PlxToolbarModule,
    PlxTooltipModule,
    PlxViewModule,
    PlxButtonsModule,
    PlxDragHandleModule,
    PlxTableModule,
    PlxChartsModule,
    PlxTabsetModule,
    // AI组件
    PlxAiGeneratorModule,
    PlxBubbleModule,
    PlxChatModule,
    PlxContentModule,
    PlxFrameNativeModule,
    PlxSenderModule,
    PlxTaskModule,
    PlxThinkingModule,
    PlxWorkflowModule,
  ]
})
export class PaletXModule { }