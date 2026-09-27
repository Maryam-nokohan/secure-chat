import { Drawer } from "@/shared/ui/Drawer";
import { useMediaQuery } from "@/shared/hooks/useMediaQuery";
import { Sidebar } from "../components/sidebar/Sidebar";
import { SettingsPanel } from "../components/settings/SettingsPanel";
import { ChatWindow } from "../components/window/ChatWindow";
import { CreateRoomModal } from "../components/modals/CreateRoomModal";
import { JoinRoomModal } from "../components/modals/JoinRoomModal";
import { AddContactModal } from "../components/modals/AddContactModal";
import { useSwipeSidebar } from "../hooks/useSwipeSidebar";
import { useChat } from "../state/ChatProvider";
import { useEffect } from "react";

export function ChatPage() {
  const { sidebarOpen, setSidebarOpen, loading, roomModal, setRoomModal } =
    useChat();
  const isDesktop = useMediaQuery("(min-width: 900px)");

  useEffect(() => {
    setSidebarOpen(isDesktop);
  }, [isDesktop]);

  useSwipeSidebar(
    sidebarOpen,
    () => setSidebarOpen(true),
    () => setSidebarOpen(false),
  );

  if (loading) {
    return (
      <div
        className="flex h-dvh items-center justify-center"
        style={{ background: "var(--surface-bg)" }}
      >
        <div
          className="h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
          style={{
            borderColor: "var(--theme-accent)",
            borderTopColor: "transparent",
          }}
        />
      </div>
    );
  }

  return (
    <div
      className="flex h-dvh w-full overflow-hidden"
      style={{ background: "var(--surface-bg)" }}
    >
      <Drawer
        open={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        static={isDesktop}
        width={300}
      >
        <Sidebar />
      </Drawer>
      <ChatWindow />
      <SettingsPanel />

      {/* Rendered here, not inside the Drawer: the Drawer's translateX()
          establishes a new containing block for position:fixed children,
          which was clipping these modals to the sidebar's width. */}
      <CreateRoomModal
        open={roomModal === "create"}
        onClose={() => setRoomModal("none")}
      />
      <JoinRoomModal
        open={roomModal === "join"}
        onClose={() => setRoomModal("none")}
      />
      <AddContactModal
        open={roomModal === "contact"}
        onClose={() => setRoomModal("none")}
      />
    </div>
  );
}
