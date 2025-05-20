<template>
  <div ref="editorContainer" class="editor-container"></div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import loader from '@monaco-editor/loader';

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  }
});
const emit = defineEmits(['update:modelValue']);

const editorContainer = ref(null);
let editorInstance;
let monacoRef;

onMounted(async () => {
  monacoRef = await loader.init();

  editorInstance = monacoRef.editor.create(editorContainer.value, {
    value: formatJson(props.modelValue),
    language: 'json',
    theme: 'vs-dark',
    automaticLayout: true,
  });

  // Format on load
  await formatDocument();

  editorInstance.onDidChangeModelContent(() => {
    const newValue = editorInstance.getValue();
    if (newValue !== props.modelValue) {
      emit('update:modelValue', newValue);
    }
  });
});

watch(() => props.modelValue, async (newValue) => {
  if (editorInstance && editorInstance.getValue() !== newValue) {
    editorInstance.setValue(formatJson(newValue));
    await formatDocument();
  }
});

onBeforeUnmount(() => {
  if (editorInstance) {
    editorInstance.dispose();
  }
});

function formatJson(str) {
  try {
    return JSON.stringify(JSON.parse(str), null, 2);
  } catch {
    return str; // Return as-is if invalid JSON
  }
}

async function formatDocument() {
  if (!monacoRef || !editorInstance) return;

  await monacoRef.languages.json.jsonDefaults.setDiagnosticsOptions({
    validate: true,
    allowComments: true
  });

  await editorInstance.getAction('editor.action.formatDocument')?.run();
}
</script>

<style scoped>
.editor-container {
  height: 400px;
  width: 100%;
  border: 1px solid #ccc;
}
</style>
