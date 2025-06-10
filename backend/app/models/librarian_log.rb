
class LibrarianLog
  include Mongoid::Document
  include Mongoid::Timestamps

  field :action, type: String

  field :admin, type: Hash     
  field :librarian, type: Hash 

  field :before_update, type: Hash, default: nil
  field :after_update, type: Hash, default: nil
end
