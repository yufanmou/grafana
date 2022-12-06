export type PluginRecipe = {
  id: string;
  name: string;
  meta: PluginRecipeMeta;
  steps: PluginRecipeStep[];
};

export type PluginRecipeMeta = {
  logo?: string;
  summary?: string;
  description?: string;
};

export type Screenshot = { name: string; url: string };

// Step - General
// --------------
export type PluginRecipeStep = {
  action: 'install-plugin' | 'display-info' | 'setup-dashboard' | 'prompt' | 'setup-alerts' | 'install-agent';

  // Meta information about the step (Optional)
  meta?: PluginRecipeStepMeta;

  // Optional information about the status of this recipe step
  status?: {
    status: string;
    statusMessage: string;
  };
};

export type PluginRecipeStepMeta = {
  name?: string;
  description?: string;
};

// Step - Install Plugin
// ---------------------
export type PluginRecipeInstallPluginStep = PluginRecipeStep & {
  meta?: PluginRecipeInstallPluginMeta;
};

export type PluginRecipeInstallPluginMeta = PluginRecipeStepMeta & {
  plugin: {
    id: string;
    version: string;
  };
};

// Step - Prompt
// --------------
export type PluginRecipePromptStep = PluginRecipeStep & {
  meta?: PluginRecipePromptMeta;
};

export type PluginRecipePromptMeta = PluginRecipeStepMeta & {
  prompts: PluginRecipePrompt[];
};

export type PluginRecipePrompt = {
  label: string;
  description: string;
  type: string;
  placeholder: string;
  defaultValue: string;
  options: Array<{
    name: string;
    value: string;
  }>;
};

// Step - Instruction
// -----------------
export type PluginRecipeInstructionStep = PluginRecipeStep & {
  meta?: PluginRecipeInstructionMeta;
};

export type PluginRecipeInstructionMeta = PluginRecipeStepMeta & {
  instructionText: string; // The text of the instructions as Markdown
  instructionTestURL: string; // The URL to run the health-check against
  instructionTestExpectedHttpResponse: string; // The expected healthy HTTP response code (e.g. "200")
};

// Step - Dashboard
// -----------------
export type PluginRecipeSetupDashboardStep = PluginRecipeStep & {
  meta?: PluginRecipeSetupDashboardStepMeta;
};

export type PluginRecipeSetupDashboardStepMeta = PluginRecipeStepMeta & {
  screenshots: Screenshot[];
};

// Step - Alerts
// -----------------

export type RecipeAlertRule = {
  namespace: string;
  group: string;
  name: string;
  summary: string;
};

export type PluginRecipeSetupAlertsStep = PluginRecipeStep & {
  alerts: RecipeAlertRule[];
};

// Step - Agent
// -----------------

export type RecipeMetric = {
  name: string;
  type?: string;
  description?: string;
};

export type PluginRecipeInstallAgentStep = PluginRecipeStep & {
  metrics: RecipeMetric[];
};

export function isSetupDashboardStep(step: PluginRecipeStep): step is PluginRecipeSetupDashboardStep {
  return step.action === 'setup-dashboard';
}

export function isInstrucitonStep(step: PluginRecipeStep): step is PluginRecipeInstructionStep {
  return step.action === 'display-info';
}

export function isPromptStep(step: PluginRecipeStep): step is PluginRecipePromptStep {
  return step.action === 'prompt';
}

export function isInstallPluginStep(step: PluginRecipeStep): step is PluginRecipeInstallPluginStep {
  return step.action === 'install-plugin';
}

export function isSetupAlertsStep(step: PluginRecipeStep): step is PluginRecipeSetupAlertsStep {
  return step.action === 'setup-alerts';
}

export function isInstallAgentStep(step: PluginRecipeStep): step is PluginRecipeInstallAgentStep {
  return step.action === 'install-agent';
}
