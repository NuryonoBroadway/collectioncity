package collection_core

type _Config struct {
	grpc      _GRPC_Config
	pubsub    _PubSub_Config
	firestore _Firestore_Config
}

type _GRPC_Config struct {
	address string
}

type _PubSub_Config struct {
	project_id      string
	project_name    string
	credential_file string
	topic           string
}

type _Firestore_Config struct {
	root_collection_id string
	root_document_id   string
}

// Config Options

type ConfigOption func(cfg *_Config) error

func WithConfig_GRPC(address string) ConfigOption {
	return func(cfg *_Config) error {
		if len(address) <= 0 {
			return Error("collection_core_client@config: grpc address is empty")
		}

		cfg.grpc.address = address

		return nil
	}
}

func WithConfig_PubSub(project_id string, project_name string, credential_file string, topic string) ConfigOption {
	return func(cfg *_Config) error {
		if len(project_id) <= 0 {
			return Error("collection_core_client@config: project id is empty")
		}

		if len(project_name) <= 0 {
			return Error("collection_core_client@config: project name is empty")
		}

		if len(credential_file) <= 0 {
			return Error("collection_core_client@config: credential file is empty")
		}

		if len(topic) <= 0 {
			return Error("collection_core_client@config: topic is empty")
		}

		cfg.pubsub = _PubSub_Config{
			project_id:      project_id,
			project_name:    project_name,
			credential_file: credential_file,
			topic:           topic,
		}

		return nil
	}
}

func WithConfig_Firestore(root_collection_id string, root_document_id string) ConfigOption {
	return func(cfg *_Config) error {
		if len(root_collection_id) <= 0 {
			return Error("collection_core_client@config: root collection id is empty")
		}

		if len(root_document_id) <= 0 {
			return Error("collection_core_client@config: root document id is empty")
		}

		cfg.firestore = _Firestore_Config{
			root_collection_id: root_collection_id,
			root_document_id:   root_document_id,
		}

		return nil
	}
}
