import { useState } from "react";
import { IconButton } from "@/shared/ui/IconButton";
import { useChat } from "../../state/ChatProvider";
import { RoomInfoModal } from "../modals/RoomInfoModal";
import { Composer } from "../composer/Composer";
import { EmptyState } from "./EmptyState";
import { MessageList } from "./MessageList";
import { Avatar } from "@/shared/ui/Avatar";

export function ChatWindow() {
  const {
    activeRoomId,
    activeRoomProfile,
    messagesByRoom,
    onlineUsers,
    rooms,
    me,
    sidebarOpen,
    setSidebarOpen,
  } = useChat();

  const [infoOpen, setInfoOpen] = useState(false);

  if (!activeRoomId || !activeRoomProfile) {
    return (
      <div className="flex h-full flex-1 flex-col">
        <MobileTopBar onMenu={() => setSidebarOpen(true)} title="Hovar" />
        <EmptyState />
      </div>
    );
  }

  const roomSummary = rooms.find((r) => r.id === activeRoomId);
  const isDirect = roomSummary?.is_direct ?? false;

  const otherMember = isDirect
    ? activeRoomProfile.members.find((m) => m.id !== me?.id)
    : undefined;

  const displayName = isDirect
    ? (otherMember?.username ?? activeRoomProfile.name)
    : activeRoomProfile.name;

  const otherOnline = otherMember ? !!onlineUsers[otherMember.id] : false;

  const onlineCount = activeRoomProfile.members.filter(
    (m) => onlineUsers[m.id],
  ).length;

  const messages = messagesByRoom[activeRoomId] ?? [];

  return (
    <div className="flex h-full min-w-0 flex-1 flex-col">
      <div
        className="flex shrink-0 items-center gap-2 px-3 py-3 sm:px-5"
        style={{
          borderBottom: "1px solid var(--surface-border)",
          background: "var(--surface-elevated)",
        }}
      >
        <button
          type="button"
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="mr-1"
          aria-label={sidebarOpen ? "Hide sidebar" : "Show sidebar"}
        >
          <MenuIcon />
        </button>
        <button
          type="button"
          onClick={() => setInfoOpen(true)}
          className="flex min-w-0 flex-1 items-center gap-2.5 text-left"
        >
          {isDirect ? (
            <Avatar username={displayName} size={36} online={otherOnline} />
          ) : (
            <span
              className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl font-display text-sm font-medium"
              style={{
                background:
                  "linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))",
                color: "var(--theme-on-accent)",
              }}
            >
              #
            </span>
          )}

          <div className="min-w-0">
            <p
              className="truncate font-display text-[15px] font-medium"
              style={{ color: "var(--text-primary)" }}
            >
              {displayName}
            </p>

            <p
              className="truncate text-xs"
              style={{
                color: otherOnline ? "#10b981" : "var(--text-tertiary)",
              }}
            >
              {isDirect
                ? otherOnline
                  ? "Online"
                  : "Offline"
                : `${activeRoomProfile.members.length} members${
                    onlineCount > 0 ? ` · ${onlineCount} online` : ""
                  }`}
            </p>
          </div>
        </button>

        <IconButton
          onClick={() => setInfoOpen(true)}
          aria-label={isDirect ? "Contact info" : "Room info"}
          title={isDirect ? "Contact info" : "Room info"}
        >
          <InfoIcon />
        </IconButton>
      </div>

      <MessageList messages={messages} />

      <Composer />

      <RoomInfoModal open={infoOpen} onClose={() => setInfoOpen(false)} />
    </div>
  );
}

function MobileTopBar({
  onMenu,
  title,
}: {
  onMenu: () => void;
  title: string;
}) {
  return (
    <div
      className="flex shrink-0 items-center gap-2 px-3 py-3 sm:hidden"
      style={{ borderBottom: "1px solid var(--surface-border)" }}
    >
      <button type="button" onClick={onMenu} aria-label="Open menu">
        <MenuIcon />
      </button>

      <span
        className="font-display text-[15px]"
        style={{ color: "var(--text-primary)" }}
      >
        {title}
      </span>
    </div>
  );
}

function MenuIcon() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="var(--text-primary)"
      strokeWidth="2"
      strokeLinecap="round"
    >
      <path d="M3 6h18M3 12h18M3 18h18" />
    </svg>
  );
}

function InfoIcon() {
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="10" />
      <path d="M12 16v-4M12 8h.01" />
    </svg>
  );
}
