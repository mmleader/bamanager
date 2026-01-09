<script>
  import { onMount } from 'svelte';
  import BrowserCard from './components/BrowserCard.svelte';
  import BrowserModal from './components/BrowserModal.svelte';
  import { 
    ListInstances, 
    AddInstance, 
    UpdateInstance, 
    DeleteInstance, 
    StartInstance, 
    StopInstance 
  } from '../wailsjs/go/main/App';

  let instances = [];
  let showModal = false;
  let editingInstance = null;

  async function loadInstances() {
    try {
      instances = await ListInstances() || [];
    } catch (err) {
      console.error("加载列表失败:", err);
    }
  }

  async function handleSave(event) {
    const data = event.detail;
    try {
      if (data.id) {
        // 更新逻辑：需要构造一个完整的 BrowserInstance 对象
        const instToUpdate = instances.find(i => i.id === data.id);
        const updated = {
            ...instToUpdate,
            name: data.name,
            path: data.path,
            userDataDir: data.userDataDir,
            args: data.args
        };
        await UpdateInstance(updated);
      } else {
        await AddInstance(data.name, data.path, data.userDataDir, data.args);
      }
      showModal = false;
      await loadInstances();
    } catch (err) {
      alert("保存失败: " + err);
    }
  }

  async function handleStart(id) {
    try {
      await StartInstance(id);
      await loadInstances();
    } catch (err) {
      alert("启动失败: " + err);
    }
  }

  async function handleStop(id) {
    try {
      await StopInstance(id);
      await loadInstances();
    } catch (err) {
      alert("停止失败: " + err);
    }
  }

  async function handleDelete(id) {
    if (!confirm("确定要删除此指纹配置吗？")) return;
    try {
      await DeleteInstance(id);
      await loadInstances();
    } catch (err) {
      alert("删除失败: " + err);
    }
  }

  function openAddModal() {
    editingInstance = null;
    showModal = true;
  }

  function openEditModal(instance) {
    editingInstance = JSON.parse(JSON.stringify(instance));
    showModal = true;
  }

  function openDuplicateModal(instance) {
    const copy = JSON.parse(JSON.stringify(instance));
    copy.id = null; // Clear ID to treat as new
    copy.name = copy.name + " (副本)";
    editingInstance = copy;
    showModal = true;
  }

  onMount(() => {
    loadInstances();
    const interval = setInterval(loadInstances, 3000); // 定时同步状态
    return () => clearInterval(interval);
  });
</script>

<main class="container">
  <header>
    <h1>指纹浏览器管理器</h1>
    <button class="btn btn-primary" on:click={openAddModal}>
      <svg style="margin-right: 0.5rem;" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
      添加实例
    </button>
  </header>

  <div class="grid">
    {#each instances as instance (instance.id)}
      <BrowserCard 
        {instance} 
        onStart={handleStart}
        onStop={handleStop}
        onEdit={openEditModal}
        onDuplicate={openDuplicateModal}
        onDelete={handleDelete}
      />
    {:else}
      <div style="grid-column: 1 / -1; text-align: center; padding: 4rem; background: white; border-radius: 0.75rem; border: 1px dashed var(--border-color); color: var(--text-muted);">
        还没有指纹例，点击上方“添加实例”开始吧。
      </div>
    {/each}
  </div>

  <BrowserModal 
    show={showModal} 
    instance={editingInstance} 
    on:save={handleSave} 
    on:close={() => showModal = false} 
  />
</main>

<style>
  /* 已经在 style.css 中定义了大部分基础样式 */
</style>
