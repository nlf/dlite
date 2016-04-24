XHYVE_VERSION := $(shell cd vendor/xhyve; git describe --abbrev=6 --dirty --always --tags)

VMM_SRC := \
	vendor/xhyve/src/vmm/x86.c \
	vendor/xhyve/src/vmm/vmm.c \
	vendor/xhyve/src/vmm/vmm_host.c \
	vendor/xhyve/src/vmm/vmm_mem.c \
	vendor/xhyve/src/vmm/vmm_lapic.c \
	vendor/xhyve/src/vmm/vmm_instruction_emul.c \
	vendor/xhyve/src/vmm/vmm_ioport.c \
	vendor/xhyve/src/vmm/vmm_callout.c \
	vendor/xhyve/src/vmm/vmm_stat.c \
	vendor/xhyve/src/vmm/vmm_util.c \
	vendor/xhyve/src/vmm/vmm_api.c \
	vendor/xhyve/src/vmm/intel/vmx.c \
	vendor/xhyve/src/vmm/intel/vmx_msr.c \
	vendor/xhyve/src/vmm/intel/vmcs.c \
	vendor/xhyve/src/vmm/io/vatpic.c \
	vendor/xhyve/src/vmm/io/vatpit.c \
	vendor/xhyve/src/vmm/io/vhpet.c \
	vendor/xhyve/src/vmm/io/vioapic.c \
	vendor/xhyve/src/vmm/io/vlapic.c \
	vendor/xhyve/src/vmm/io/vpmtmr.c \
	vendor/xhyve/src/vmm/io/vrtc.c

XHYVE_SRC := \
	vendor/xhyve/src/acpitbl.c \
	vendor/xhyve/src/atkbdc.c \
	vendor/xhyve/src/block_if.c \
	vendor/xhyve/src/consport.c \
	vendor/xhyve/src/dbgport.c \
	vendor/xhyve/src/inout.c \
	vendor/xhyve/src/ioapic.c \
	vendor/xhyve/src/md5c.c \
	vendor/xhyve/src/mem.c \
	vendor/xhyve/src/mevent.c \
	vendor/xhyve/src/mptbl.c \
	vendor/xhyve/src/pci_ahci.c \
	vendor/xhyve/src/pci_emul.c \
	vendor/xhyve/src/pci_hostbridge.c \
	vendor/xhyve/src/pci_irq.c \
	vendor/xhyve/src/pci_lpc.c \
	vendor/xhyve/src/pci_uart.c \
	vendor/xhyve/src/pci_virtio_block.c \
	vendor/xhyve/src/pci_virtio_net_tap.c \
	vendor/xhyve/src/pci_virtio_net_vmnet.c \
	vendor/xhyve/src/pci_virtio_rnd.c \
	vendor/xhyve/src/pm.c \
	vendor/xhyve/src/post.c \
	vendor/xhyve/src/rtc.c \
	vendor/xhyve/src/smbiostbl.c \
	vendor/xhyve/src/task_switch.c \
	vendor/xhyve/src/uart_emul.c \
	vendor/xhyve/src/xhyve.c \
	vendor/xhyve/src/virtio.c \
	vendor/xhyve/src/xmsr.c

FIRMWARE_SRC := \
	vendor/xhyve/src/firmware/kexec.c \
	vendor/xhyve/src/firmware/fbsd.c

SRC := \
	$(VMM_SRC) \
	$(XHYVE_SRC) \
	$(FIRMWARE_SRC)
