import React, { memo } from "react";
import { View, Text, StyleSheet } from "react-native";

interface Message {
  chat_id: number;
  content: string;
  created_at: string;
  user: {
    name: string;
  };
  user_id: number | null;
}

interface MessageItemProps {
  item: Message;
  isOwnMessage: boolean;
}

const MessageItem: React.FC<MessageItemProps> = ({ item, isOwnMessage }) => {
  return isOwnMessage ? (
    <View style={styles.ownMessage}>
      <Text>{item.content}</Text>
    </View>
  ) : (
    <View style={styles.someoneMessage}>
      <Text>
        {item.user.name}: {item.content}
      </Text>
    </View>
  );
};

export default memo(MessageItem);

const styles = StyleSheet.create({
  ownMessage: {
    alignSelf: "flex-end",
    backgroundColor: "#9BFFAE",
    padding: 10,
    borderRadius: 10,
    marginBottom: 5,
    maxWidth: "75%",
  },
  someoneMessage: {
    alignSelf: "flex-start",
    backgroundColor: "#B2B4B2",
    padding: 10,
    borderRadius: 10,
    marginBottom: 5,
    maxWidth: "75%",
  },
});
