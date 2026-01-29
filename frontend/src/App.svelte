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
    StopInstance,
    GetConfig,
    SetMinimizeToTray
  } from '../wailsjs/go/main/App';

  let instances = [];
  let showModal = false;

  let editingInstance = null;
  let minimizeToTray = false;

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

  onMount(async () => {
    try {
      const config = await GetConfig();
      if (config) {
        minimizeToTray = config.minimize_to_tray;
      }
    } catch (e) {
      console.error("加载配置失败", e);
    }
    loadInstances();
    const interval = setInterval(loadInstances, 3000); // 定时同步状态
    return () => clearInterval(interval);
  });
</script>

<main class="container">
  <header>
    <h1>指纹浏览器管理器</h1>
    <div style="display: flex; align-items: center; gap: 1rem;">
      <label class="toggle-switch">
        <input type="checkbox" bind:checked={minimizeToTray} on:change={(e) => SetMinimizeToTray(e.target.checked)} />
        <span class="slider"></span>
        <span class="label-text">关闭时最小化到托盘</span>
      </label>
      <button class="btn btn-primary" on:click={openAddModal}>
        <svg style="margin-right: 0.5rem;" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        添加实例
      </button>
    </div>
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
  
  /* Toggle Switch Styles */
  .toggle-switch {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    user-select: none;
    font-size: 0.9rem;
    color: var(--text-main);
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
    position: absolute; /* Remove from flow */
  }

  .slider {
    position: relative;
    width: 44px;
    height: 24px;
    background-color: #cbd5e1; /* slate-300 */
    border-radius: 9999px;
    transition: .3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 20px;
    width: 20px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    transition: .3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 50%;
    box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
  }

  input:checked + .slider {
    background-color: var(--primary-color);
  }

  input:checked + .slider:before {
    transform: translateX(20px);
  }

  .label-text {
    font-weight: 500;
  }
</style>
